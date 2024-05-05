package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/bradleyjkemp/sigma-go"
	"gopkg.in/yaml.v3"
)

var (
	app = kingpin.New("velosigma",
		"A tool for manipulating sigma files.")

	compile_cmd             = app.Command("compile", "Compile all the rules into one rule.")
	output                  = compile_cmd.Flag("output", "File to write the artifact bundle to").String()
	yaml_output             = compile_cmd.Flag("yaml", "File to write the artifact yaml to").String()
	rejects_output          = compile_cmd.Flag("rejects", "File to write the rejected rules to").String()
	ignore_previous_rejects = compile_cmd.Flag("ignore_previous_rejects", "Read the rejects file and ignore any previously known rejected rules").Bool()
	config                  = compile_cmd.Flag("config", "Config file to use").Required().ExistingFile()
	level_regex_str         = compile_cmd.Flag("level_regex", "A regex to select rule Levels").Default(".").String()

	debug = app.Flag("debug", "Print more details").Bool()

	command_handlers []CommandHandler

	allowed_additional_fields = []string{"details", "vql", "vql_args"}
)

type CommandHandler func(command string) bool

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	for _, handler := range command_handlers {
		if handler(command) {
			break
		}
	}
}

func (self *CompilerContext) CompileDirs() error {
	for _, compile_dir := range self.config_obj.RuleDirectories {
		fmt.Printf("Scanning directory %v for rules\n", compile_dir)
		err := filepath.WalkDir(compile_dir,
			func(path string, d fs.DirEntry, err error) error {
				if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
					return nil
				}

				fd, err := os.Open(path)
				if err != nil {
					return err
				}

				data, err := ioutil.ReadAll(fd)
				if err != nil {
					return err
				}

				return self.CompileRule(string(data), path)
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *CompilerContext) CompileRule(rule_yaml, path string) error {
	rule, err := sigma.ParseRule([]byte(rule_yaml))
	if err != nil {
		self.addError(err.Error(), path)
		if !self.shouldSuppressError(err.Error(), path) {
			fmt.Printf("While compiling rule %v: %v\n", path, err)
		}
		return nil
	}

	if rule.Level == "" {
		rule.Level = "default"
	}

	if !self.level_regex.MatchString(rule.Level) {
		return nil
	}

	additional_fields := make(map[string]interface{})
	for _, f := range allowed_additional_fields {
		v := rule.AdditionalFields[f]
		if v != nil {
			additional_fields[f] = v
		}
	}

	new_rule := sigma.Rule{
		Title:            rule.Title,
		Author:           rule.Author,
		Level:            rule.Level,
		Status:           rule.Status,
		Logsource:        rule.Logsource,
		Detection:        rule.Detection,
		AdditionalFields: additional_fields,
	}

	if self.config_obj.BaseReferenceURL != "" {
		new_rule.References = []string{
			self.config_obj.BaseReferenceURL + path}
	}

	// Record all the rules we added
	self.total_visited_rules++

	// Skip errored rules.
	logsource, err := self.normalize_logsource(&new_rule, path)
	if err != nil {
		self.addError(err.Error(), path)
		return nil
	}

	err = self.walk_fields(&new_rule, path, logsource)
	if err != nil {
		self.addError(err.Error(), path)
		return nil
	}

	err = self.check_condition(&new_rule)
	if err != nil {
		self.addError(err.Error(), path)
		return nil
	}

	yamlEncoder := yaml.NewEncoder(&self.rules)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(new_rule)
	if err != nil {
		self.addError(err.Error(), path)
		return nil
	}

	DebugPrint("Processing %v\n", path)
	self.rules.Write([]byte("\n---\n"))
	self.incLogSource(logsource)

	// Only write the original_rules we actually added -
	// rejected rules will not be added to the zip file.
	self.original_rules.Write([]byte(rule_yaml))
	self.original_rules.Write([]byte("\n---\n"))

	return nil
}

func doCompile() (err error) {
	if *yaml_output == "" && *output == "" {
		return errors.New("Must provide either --output or --yaml")
	}

	level_regex, err := regexp.Compile(*level_regex_str)
	if err != nil {
		return fmt.Errorf("Level Regex invalid: %w", err)
	}

	context := NewCompilerContext()
	context.level_regex = level_regex

	err = context.LoadConfig(*config)
	if err != nil {
		return fmt.Errorf("Reading Config: %w", err)
	}

	if *ignore_previous_rejects && *rejects_output != "" {
		err := context.LoadRejectSupporessions(*rejects_output)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Reading Rejects: %w", err)
		}
	}

	defer func() {
		fmt.Printf("\nGenerated rules with level %v into %v\n",
			*level_regex_str, *output)
		stats := context.Stats()

		fmt.Printf("\nTotal rules added: %v from %v visited files and %v rejected rules\n",
			stats.TotalRules, stats.TotalVisitedRules, stats.TotalErrors)

		if stats.TotalUnhandledErrors > 0 {
			err = fmt.Errorf("Unhandled errors %v", stats.TotalUnhandledErrors)
		}
	}()

	err = context.CompileDirs()
	if err != nil {
		return fmt.Errorf("Listing directory: %w", err)
	}

	if *output != "" {
		// Write the sigma file in the output directory.
		out_fd, err := os.OpenFile(*output,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("Creating output: %w", err)
		}
		defer out_fd.Close()

		zip := zip.NewWriter(out_fd)
		defer zip.Close()

		err = context.WriteArtifact(zip)
		if err != nil {
			return fmt.Errorf("WriteArtifact: %w", err)
		}
	}

	if *yaml_output != "" {
		out_fd, err := os.OpenFile(*yaml_output,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("Creating yaml output: %w", err)
		}
		defer out_fd.Close()

		artifact, err := context.GetArtifact()
		if err != nil {
			return fmt.Errorf("GetArtifact: %w", err)
		}

		out_fd.Write([]byte(artifact))
	}

	if *rejects_output != "" {
		out_fd, err := os.OpenFile(*rejects_output,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("Creating yaml output: %w", err)
		}
		defer out_fd.Close()

		out_fd.Write(MustMarshalIndent(context.GetRejected()))
	}

	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case compile_cmd.FullCommand():
			err := doCompile()
			kingpin.FatalIfError(err, "Compiling artifact")

		default:
			return false
		}
		return true
	})
}
