package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Velocidex/sigma-go"
	kingpin "github.com/alecthomas/kingpin/v2"
	"gopkg.in/yaml.v3"
)

var (
	app = kingpin.New("velosigma",
		"A tool for manipulating sigma files.")

	compile_cmd             = app.Command("compile", "Compile all the rules into one rule.")
	output                  = compile_cmd.Flag("output", "File to write the artifact bundle to").String()
	yaml_output             = compile_cmd.Flag("yaml", "File to write the artifact yaml to").String()
	docs_output             = compile_cmd.Flag("docs", "File to write the documentation markdown to").String()
	rejects_output          = compile_cmd.Flag("rejects", "File to write the rejected rules to").String()
	ignore_previous_rejects = compile_cmd.Flag("ignore_previous_rejects", "Read the rejects file and ignore any previously known rejected rules").Bool()
	config                  = compile_cmd.Flag("config", "Config file to use").Required().ExistingFile()
	level_regex_str         = compile_cmd.Flag("level_regex", "A regex to select rule Levels").Default(".").String()
	rule_directory          = compile_cmd.Flag("rule_dir", "The base directory to write the rules").String()

	debug = app.Flag("debug", "Print more details").Bool()

	command_handlers []CommandHandler

	allowed_additional_fields = []string{"details", "vql", "vql_args", "enrichment"}
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

				data, err := io.ReadAll(fd)
				if err != nil {
					return err
				}

				// Allow each file to contain multiple rules.
				for _, rule := range strings.Split(string(data), "\n---\n") {
					rule_data := strings.TrimSpace(rule) + "\n"
					if rule_data == "" {
						continue
					}

					err := self.CompileRule(rule_data, path)
					if err != nil {
						fmt.Printf("Rule %v: %v\n", err, rule_data)
						return err
					}
				}

				return nil
			})
		if err != nil {
			return err
		}
	}
	return self.checkDefaultDetails()
}

func (self *CompilerContext) checkDefaultDetails() (err error) {
	if self.config_obj.DefaultDetails.Lookup != nil {
		for k, details := range self.config_obj.DefaultDetails.Lookup {
			err = self.walk_details(details, k)
			if err != nil {
				fmt.Printf("Checking Default Details: %v\n", err)
			}
		}
	}

	return err
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

	if rule.Title == "" {
		self.addError("Rule without Title", path)
		return nil
	}

	// Sometimes rules have the same title so we make the title
	// unique.
	count, pres := self.seen_rules[rule.Title]
	if pres {
		self.seen_rules[rule.Title] = count + 1
		rule.Title = fmt.Sprintf("%v %d", rule.Title, count+1)
	} else {
		self.seen_rules[rule.Title] = 1
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

	// This is the concise and reducted rule that will be hunted for.
	new_rule := sigma.Rule{
		Title:       rule.Title,
		Author:      rule.Author,
		Level:       rule.Level,
		Status:      rule.Status,
		Logsource:   rule.Logsource,
		Detection:   self.normalize_detections(rule.Detection),
		Correlation: rule.Correlation,
		//Detection:        rule.Detection,
		AdditionalFields: additional_fields,
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

	buf := &bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(buf)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(new_rule)
	if err != nil {
		self.addError(err.Error(), path)
		return nil
	}

	DebugPrint("Processing %v\n", path)
	self.rules = append(self.rules, buf.String())

	self.incLogSource(logsource)

	// Only write the original_rules we actually added -
	// rejected rules will not be added to the zip file.
	self.original_rules.Write([]byte(rule_yaml))
	self.original_rules.Write([]byte("\n---\n"))

	self.original_rules_by_path[path] = rule_yaml
	self.original_rule_obj_by_path[path] = rule
	self.rules_by_path[path] = new_rule

	return nil
}

func (self *CompilerContext) normalize_search(
	search sigma.Search) sigma.Search {

	res := sigma.Search{}

	// Specifically handle a search with no field but a set of
	// keywords. For example:
	//
	// detection:
	//   keywords_filter:
	//      - 'of='
	//
	// I think this means match any of the keywords on the event
	if len(search.EventMatchers) == 0 && len(search.Keywords) > 0 {
		new_matcher := sigma.FieldMatcher{
			Field: "Event",
		}
		for _, kw := range search.Keywords {
			new_matcher.Values = append(new_matcher.Values, kw)
		}
		res.EventMatchers = []sigma.EventMatcher{
			[]sigma.FieldMatcher{new_matcher}}
		return res
	}

	for _, field_matcher := range search.EventMatchers {
		new_field_matchers := []sigma.FieldMatcher{}

		// The following code expands strings in hex format to a OR
		// condition with either the hex value or the decimal
		// value. This is because Velociraptor's evtx parser preserves
		// native types and so will emit the value as an integer, but
		// many rules are matching against the value as a hex encoded
		// number. In that case, where it is safe, we expand the
		// condition to be a logical OR of the hex and decimal
		// versions of the number.
		for _, m := range field_matcher {
			// Some rules do not name a field ar all. In that case we
			// assume the field is called "Event" and let the log
			// source do something sensible with it.
			if m.Field == "" {
				m.Field = "Event"
			}
			// Do not touch 'all' modifiers which represent a logical
			// AND. If we expand them into an OR we can change the
			// logic.
			if !InString(m.Modifiers, "all") {
				new_values := make([]interface{}, 0, len(m.Values))
				for _, v := range m.Values {
					v_str, ok := v.(string)
					if ok {
						if strings.HasPrefix(v_str, "0x") {
							v_int, err := strconv.ParseUint(v_str, 0, 64)
							if err == nil {
								new_values = append(new_values,
									fmt.Sprintf("%d", v_int))
							}
						}
					}
					new_values = append(new_values, v)
				}
				m.Values = new_values
			}
			new_field_matchers = append(new_field_matchers, m)
		}

		res.EventMatchers = append(res.EventMatchers, new_field_matchers)
	}

	return res
}

func (self *CompilerContext) normalize_detections(
	detections sigma.Detection) sigma.Detection {

	for k, v := range detections.Searches {
		detections.Searches[k] = self.normalize_search(v)
	}
	return detections
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
		out_fd, err := CreateFile(*output)
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
		out_fd, err := CreateFile(*yaml_output)
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

	if *docs_output != "" {
		out_fd, err := CreateFile(*docs_output)
		if err != nil {
			return fmt.Errorf("Creating docs output: %w", err)
		}
		defer out_fd.Close()

		documentation, err := context.GetDocs()
		if err != nil {
			return fmt.Errorf("GetDocs: %w", err)
		}

		out_fd.Write([]byte(documentation))
	}

	if *rule_directory != "" {
		err := context.WriteRuleDir(*rule_directory)
		if err != nil {
			return fmt.Errorf("WriteRuleDir: %w", err)
		}
	}

	if *rejects_output != "" {
		out_fd, err := CreateFile(*rejects_output)
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
