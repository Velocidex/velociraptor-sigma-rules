package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Velocidex/yaml/v2"
	"github.com/bradleyjkemp/sigma-go"
)

type CompilerContext struct {
	logsources map[string]int

	// Keep track of fields seen for each log source that are not
	// already known and defined by the config file.
	missing_fields_in_logsources map[string]map[string][]string

	fields         map[string]int
	missing_fields map[string][]string

	// Invalid fields are available fields, but not for its logsource
	invalid_fields map[string][]string

	// Map between error reason and the files that were rejected for
	// it.
	errored_rules map[string][]string
	config_obj    *Config

	rules       bytes.Buffer
	level_regex *regexp.Regexp

	total_visited_rules int

	vql bytes.Buffer

	original_rules *bytes.Buffer
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		logsources:                   make(map[string]int),
		missing_fields_in_logsources: make(map[string]map[string][]string),
		errored_rules:                make(map[string][]string),
		level_regex:                  regexp.MustCompile(".*"),

		fields:         make(map[string]int),
		missing_fields: make(map[string][]string),
		invalid_fields: make(map[string][]string),
		original_rules: &bytes.Buffer{},
	}
}

func (self *CompilerContext) LoadConfig(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	return self.LoadConfigFromString(string(data))
}

func (self *CompilerContext) LoadConfigFromString(data string) error {

	config_obj := &Config{}
	err := yaml.Unmarshal([]byte(data), config_obj)
	if err != nil {
		return err
	}
	self.config_obj = config_obj

	return nil
}

func (self *CompilerContext) Resolve(source_spec string) bool {
	_, pres := self.config_obj.Sources[source_spec]
	return pres
}

func (self *CompilerContext) incLogSource(source_spec string) {
	count, _ := self.logsources[source_spec]
	count++

	self.logsources[source_spec] = count
}

func (self *CompilerContext) Stats() {
	total_rules := 0
	sources := []string{}
	for k := range self.logsources {
		sources = append(sources, k)
	}

	sort.Strings(sources)
	fmt.Printf("\nThe following log sources will be used:\n")
	for _, v := range sources {
		count, _ := self.logsources[v]

		fmt.Printf("  %v (%v rules)\n", v, count)
		total_rules += count
	}

	total_reject_rules := 0
	if len(self.errored_rules) > 0 {
		fmt.Printf("\nErrored Rules which were rejected:\n")
		for k, v := range self.errored_rules {
			fmt.Printf("  %v:\n", k)
			for _, i := range v {
				fmt.Printf("     %v\n", i)
				total_reject_rules++
			}
		}
	}

	if len(self.missing_fields_in_logsources) > 0 {
		/*
			fmt.Printf("\nRules that reference unknown fields in each log source:\n")
			for logsource, field_map := range self.missing_fields_in_logsources {
				fmt.Printf("  %v:\n", logsource)
				for field, path_list := range field_map {
					fmt.Printf("     %v:\n", field)

					for _, path := range path_list {
						fmt.Printf("         %v\n", path)
					}
				}
			}
		*/
		fmt.Printf("\nFields by Log source:\n")
		for logsource, field_map := range self.missing_fields_in_logsources {
			fmt.Printf("%v:\n", logsource)
			for field, _ := range field_map {
				fmt.Printf("  - %v\n", field)
			}
		}

	}

	fmt.Printf("\nTotal rules added: %v from %v visited files and %v rejected rules\n",
		total_rules, self.total_visited_rules, total_reject_rules)

	return

	if len(self.missing_fields) > 0 {
		fmt.Printf("\nMissing field mappings:\n")
		for k, v := range self.missing_fields {
			fmt.Printf("  %v:\n", k)
			for _, i := range v {
				fmt.Printf("     %v\n", i)
			}
		}
	}

	if len(self.invalid_fields) > 0 {
		fmt.Printf("\nInvalid field mappings:\n")
		for k, v := range self.invalid_fields {
			fmt.Printf("  %v:\n", k)
			for _, i := range v {
				fmt.Printf("     %v\n", i)
			}
		}
	}
}

func (self *CompilerContext) getSourceFromChannel(
	channel string) string {
	for source, query := range self.config_obj.Sources {
		for _, c := range query.Channel {
			if c == channel {
				return source
			}
		}
	}
	return ""
}

// Update the rule's log source to specify the logsource
func (self *CompilerContext) updateRuleLogSources(source_spec string, rule *sigma.Rule) {
	parts := strings.Split(source_spec, "/")
	if len(parts) == 3 {
		if parts[0] == "*" {
			parts[0] = ""
		}

		rule.Logsource.Category = parts[0]
		rule.Logsource.Product = parts[1]
		rule.Logsource.Service = parts[2]
	}
}

func (self *CompilerContext) check_condition(rule *sigma.Rule) error {
	// If a condition is missing make it up
	if len(rule.Detection.Conditions) == 0 {
		fields := []string{}
		for k := range rule.Detection.Searches {
			fields = append(fields, k)
		}

		rule.Detection.Conditions = append(rule.Detection.Conditions,
			sigma.Condition{
				Search: sigma.SearchIdentifier{
					Name: strings.Join(fields, " and "),
				},
			})
	}
	return nil
}

func (self *CompilerContext) guessLogSource(
	rule *sigma.Rule,
	category, product, service string) (reason string, source string) {

	// Try to find a Channel match
	for _, search := range rule.Detection.Searches {
		for _, event_matcher := range search.EventMatchers {
			for _, matcher := range event_matcher {
				if matcher.Field == "Channel" {
					for _, v := range matcher.Values {
						v_str, ok := v.(string)
						if ok {
							source := self.getSourceFromChannel(v_str)
							if source != "" {
								self.updateRuleLogSources(source, rule)
								return fmt.Sprintf(
									"Found Channel %v", v_str), source
							}
						}
					}
				}
			}
		}
	}

	return "", fmt.Sprintf("%v/%v/%v", category, product, service)
}

// Keep track of Sigma rules that refer to a field which is not
// defined in the field mapping.
func (self *CompilerContext) incMissingFieldMap(field, path string) {
	missing, _ := self.missing_fields[field]
	missing = append(missing, path)
	self.missing_fields[field] = missing
}

// Keep track of fields in sigma rules that refer to underfined fields
// for their specified log source.
func (self *CompilerContext) incMissingFieldInLogSource(
	field, logsource, path string) {

	field_map, pres := self.missing_fields_in_logsources[logsource]
	if !pres {
		field_map = make(map[string][]string)
	}
	self.missing_fields_in_logsources[logsource] = field_map

	path_list, _ := field_map[field]
	path_list = append(path_list, path)
	field_map[field] = path_list
}

var valid_modifiers = map[string]bool{
	"all":        true,
	"any":        true,
	"re":         true,
	"contains":   true,
	"endswith":   true,
	"startswith": true,
	"cidr":       true,
	"gt":         true,
	"gte":        true,
	"lt":         true,
	"lte":        true,
	"vql":        true,
}

func (self *CompilerContext) check_modifiers(
	matcher sigma.FieldMatcher, path string) error {
	for _, m := range matcher.Modifiers {
		_, pres := valid_modifiers[m]
		if !pres {
			err := fmt.Errorf("Invalid modifier %v", m)
			self.addError(err.Error(), path)
			return err
		}
	}

	return nil
}

func (self *CompilerContext) walk_fields(
	rule *sigma.Rule, path string, logsource string) (err error) {
	for _, search := range rule.Detection.Searches {
		for _, event_matcher := range search.EventMatchers {
			for _, matcher := range event_matcher {
				// Check if modifiers are valid.
				err := self.check_modifiers(matcher, path)
				if err != nil {
					return err
				}

				// Check if there is a field mapping
				_, pres := self.config_obj.FieldMappings[matcher.Field]
				if !pres {
					self.incMissingFieldMap(matcher.Field, path)
					return fmt.Errorf(
						"Missing field mapping '%v' in %v", matcher.Field, logsource)
				}

				// Just report an unknown field but do not reject the
				// rule - this is just a warning that the rule is
				// using a field on this log source which is not known
				// to belong to this log source.
				pres = slices.Contains(
					self.config_obj.Sources[logsource].Fields, matcher.Field)
				if !pres {
					self.incMissingFieldInLogSource(matcher.Field,
						logsource, path)
				}
			}
		}
	}
	return err
}

func (self *CompilerContext) normalize_logsource(
	rule *sigma.Rule, path string) (string, error) {
	source := rule.Logsource
	category := source.Category
	if category == "" {
		category = "*"
	}

	product := source.Product
	if product == "" {
		product = "*"
	}

	service := source.Service
	if service == "" {
		service = "*"
	}

	source_spec := fmt.Sprintf("%v/%v/%v", category, product, service)

	// Try to find a suitable log source based on rule inspection itself.
	reason, guessed_source_spec := self.guessLogSource(rule, category, product, service)
	if self.Resolve(guessed_source_spec) &&
		guessed_source_spec != source_spec {
		DebugPrint("** Substitute guess source %v for %v with rule %v %v\n",
			guessed_source_spec, source_spec, path, reason)
		return guessed_source_spec, nil
	}

	// Otherwise see if we have a direct log source defined.
	if self.Resolve(source_spec) {
		return source_spec, nil

	}
	// Try to guess the source spec
	DebugPrint("**** Log Source '%v' not found!\n", guessed_source_spec)
	return source_spec, fmt.Errorf(
		"Missing Source: '%v'", source_spec)
}

func (self *CompilerContext) addError(reason string, path string) {
	failed, _ := self.errored_rules[reason]
	failed = append(failed, path)
	self.errored_rules[reason] = failed
}
