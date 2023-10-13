package main

import (
	"archive/zip"
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

	compile_cmd     = app.Command("compile", "Compile all the rules into one rule.")
	compile_dirs    = compile_cmd.Arg("directory", "Directory to compile").Required().Strings()
	output          = compile_cmd.Flag("output", "File to write the artifact to").Required().String()
	config          = compile_cmd.Flag("config", "Config file to use").Required().ExistingFile()
	level_regex_str = compile_cmd.Flag("level_regex", "A regex to select rule Levels").Default(".").String()

	base_url = compile_cmd.Flag("base_url", "URL to rule").
			Default("https://github.com/Yamato-Security/hayabusa-rules/tree/main/").
			String()

	command_handlers []CommandHandler
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

func (self *CompilerContext) CompileDirs(directories []string) error {
	for _, compile_dir := range directories {
		err := filepath.WalkDir(compile_dir,
			func(path string, d fs.DirEntry, err error) error {
				if !strings.HasSuffix(path, ".yml") {
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

				rule, err := sigma.ParseRule(data)
				if err != nil {
					return nil
				}

				if !self.level_regex.MatchString(rule.Level) {
					return nil
				}

				additional_fields := make(map[string]interface{})
				details := rule.AdditionalFields["details"]
				if details != nil {
					additional_fields["details"] = details
				}

				new_rule := sigma.Rule{
					Title:            rule.Title,
					Author:           rule.Author,
					Level:            rule.Level,
					Status:           rule.Status,
					Logsource:        rule.Logsource,
					Detection:        rule.Detection,
					AdditionalFields: additional_fields,
					References:       []string{*base_url + path},
				}

				self.normalize_logsource(&new_rule, path)
				err = self.walk_fields(&new_rule, path)
				if err != nil {
					return nil
				}

				yamlEncoder := yaml.NewEncoder(&self.rules)
				yamlEncoder.SetIndent(2)
				err = yamlEncoder.Encode(new_rule)
				if err != nil {
					return nil
				}

				fmt.Printf("Processing %v\n", path)
				self.rules.Write([]byte("\n---\n"))

				return err
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func doCompile() {
	// Write the sigma file in the output directory.
	out_fd, err := os.OpenFile(*output,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	kingpin.FatalIfError(err, "Creating output")
	defer out_fd.Close()

	level_regex, err := regexp.Compile(*level_regex_str)
	kingpin.FatalIfError(err, "Level Regex invalid")

	zip := zip.NewWriter(out_fd)
	defer zip.Close()

	context := NewCompilerContext()
	context.level_regex = level_regex

	err = context.LoadConfig(*config)
	kingpin.FatalIfError(err, "Reading Config")

	defer func() {
		fmt.Printf("Generated rules with level %v into %v\n",
			*level_regex_str, *output)
		context.Stats()
	}()

	err = context.CompileDirs(*compile_dirs)
	kingpin.FatalIfError(err, "Listing directory")

	err = context.WriteArtifact(zip)
	kingpin.FatalIfError(err, "WriteArtifact")
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case compile_cmd.FullCommand():
			doCompile()

		default:
			return false
		}
		return true
	})
}
