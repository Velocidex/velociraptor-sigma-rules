package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Velocidex/ordereddict"
	"github.com/Velocidex/sigma-go"
	"github.com/Velocidex/yaml/v2"
)

type Stats struct {
	TotalUnhandledErrors int
	TotalErrors          int
	TotalRules           int
	TotalVisitedRules    int
}

type CompilerContext struct {
	logsources       map[string]int
	logsources_order []string

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

	// We load the previous run's rejected rules to only show
	// incremental failures
	ignored_rules map[string]bool
	seen_rules    map[string]int

	config_obj *Config

	rules       []string
	level_regex *regexp.Regexp

	total_visited_rules int

	vql bytes.Buffer

	original_rules *bytes.Buffer

	event_resolver *EventResolver

	imported_configs []*Config

	// Keep the original rules as raw strings to include any comments
	// in the yaml which would be lost on decode/encode round trip.
	original_rules_by_path    map[string]string
	original_rule_obj_by_path map[string]sigma.Rule
	rules_by_path             map[string]sigma.Rule

	completed_artifact string
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		logsources:                   make(map[string]int),
		missing_fields_in_logsources: make(map[string]map[string][]string),
		seen_rules:                   make(map[string]int),
		errored_rules:                make(map[string][]string),
		ignored_rules:                make(map[string]bool),
		original_rules_by_path:       make(map[string]string),
		original_rule_obj_by_path:    make(map[string]sigma.Rule),
		rules_by_path:                make(map[string]sigma.Rule),
		level_regex:                  regexp.MustCompile(".*"),

		fields:         make(map[string]int),
		missing_fields: make(map[string][]string),
		invalid_fields: make(map[string][]string),
		original_rules: &bytes.Buffer{},
		event_resolver: &EventResolver{},
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

	config_obj := NewConfig()
	err := yaml.Unmarshal([]byte(data), config_obj)
	if err != nil {
		return err
	}
	self.config_obj = config_obj

	if config_obj.EventResolver != "" {
		err = self.event_resolver.Load(config_obj.EventResolver)
		if err != nil {
			return err
		}
	}

	self.event_resolver.config_obj = config_obj

	for _, path := range config_obj.ImportConfigs {
		fd, err := os.Open(path)
		if err != nil {
			return err
		}

		secondary_data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		config_obj = NewConfig()
		err = yaml.Unmarshal(secondary_data, config_obj)
		if err != nil {
			return err
		}

		self.imported_configs = append(self.imported_configs, config_obj)
		self.config_obj.mergeConfig(config_obj)
	}
	self.config_obj.mergeConfig(self.config_obj)

	return nil
}

func (self *CompilerContext) LoadRejectSupporessions(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	rejects := &Rejected{}
	err = json.Unmarshal(data, rejects)
	if err != nil {
		return err
	}

	for _, item := range rejects.Rejects {
		self.ignored_rules[item.String()] = true
	}

	return nil
}

// Get the source_spec from either this config or imported_configs
func (self *CompilerContext) Resolve(source_spec string) bool {
	_, pres := self.config_obj.sources.Get(source_spec)
	return pres
}

func (self *CompilerContext) incLogSource(source_spec string) {
	count, pres := self.logsources[source_spec]
	if !pres {
		self.logsources_order = append(self.logsources_order, source_spec)
	}
	count++

	self.logsources[source_spec] = count
}

func (self *CompilerContext) shouldSuppressError(err_msg string, path string) bool {
	for _, m := range self.config_obj.BadFieldMappings {
		if strings.Contains(err_msg, m) {
			return true
		}
	}

	for _, m := range self.config_obj.BadSources {
		if strings.Contains(err_msg, m) {
			return true
		}
	}

	key := path + err_msg
	_, pres := self.ignored_rules[key]
	if pres {
		return true
	}

	return false
}

func (self *CompilerContext) Stats() Stats {
	result := Stats{
		TotalVisitedRules: self.total_visited_rules,
	}

	fmt.Printf("\nThe following log sources will be used:\n")
	for _, v := range self.logsources_order {
		count, _ := self.logsources[v]

		fmt.Printf("  %v (%v rules)\n", v, count)
		result.TotalRules += count
	}

	if len(self.errored_rules) > 0 {
		fmt.Printf("\nErrored Rules which were rejected:\n")
		for k, v := range self.errored_rules {
			messages := []string{}
			for _, path := range v {
				if self.shouldSuppressError(k, path) {
					result.TotalErrors += len(v)
					continue
				}

				messages = append(messages, fmt.Sprintf("     %v", path))
				result.TotalUnhandledErrors++
			}

			if len(messages) > 0 {
				fmt.Printf("  %v:\n", k)
				for _, m := range messages {
					fmt.Println(m)
				}
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

	return result

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

	return result
}

func (self *CompilerContext) getSourceFromChannel(
	channel string) string {
	for _, source := range self.config_obj.sources.Keys() {
		query_any, pres := self.config_obj.sources.Get(source)
		if !pres {
			continue
		}

		query, ok := query_any.(Query)
		if !ok {
			continue
		}

		for _, c := range query.Channel {
			if c == channel {
				return source
			}
		}
	}
	return ""
}

// Update the rule's log source to specify the logsource
func updateRuleLogSources(source_spec string, log_source *sigma.Logsource) {
	parts := strings.Split(source_spec, "/")
	if len(parts) == 3 {
		if parts[0] == "*" {
			parts[0] = ""
		}
		if parts[1] == "*" {
			parts[1] = ""
		}
		if parts[2] == "*" {
			parts[2] = ""
		}
		log_source.Category = parts[0]
		log_source.Product = parts[1]
		log_source.Service = parts[2]
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

	if rule.Detection.Timeframe != "" {
		return fmt.Errorf("In rule %v: Timeframe detections not supported",
			rule.Title)
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
								updateRuleLogSources(source, &rule.Logsource)
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
	"all":          true,
	"any":          true,
	"re":           true,
	"contains":     true,
	"endswith":     true,
	"startswith":   true,
	"cidr":         true,
	"gt":           true,
	"gte":          true,
	"lt":           true,
	"lte":          true,
	"vql":          true,
	"base64":       true,
	"base64offset": true,
	"windash":      true,
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

var expandRegEx = regexp.MustCompile(`%([A-Z.a-z_0-9]+)(\[([0-9]+)\])?%`)

func (self *CompilerContext) walk_details(
	details string, path string, logsource string) (err error) {

	for _, match := range expandRegEx.FindAllStringSubmatch(details, -1) {
		variable := match[1]

		if variable == "UserContext" {
			DlvBreak()
		}

		// Check if there is a field mapping
		if !self.event_resolver.CheckFieldMapping(variable) {
			self.incMissingFieldMap(variable, path)
			return fmt.Errorf(
				"Missing field mapping '%v' in %v", variable, logsource)
		}
	}

	return nil
}

func (self *CompilerContext) walk_fields(
	rule *sigma.Rule, path string, logsource string) (err error) {

	// This needs to be stable so the errors are consistent
	for _, search := range ValuesInOrder(rule.Detection.Searches) {
		for _, event_matcher := range search.EventMatchers {
			for _, matcher := range event_matcher {
				// Check if modifiers are valid.
				err := self.check_modifiers(matcher, path)
				if err != nil {
					return err
				}

				// Check if there is a field mapping
				if !self.event_resolver.CheckFieldMapping(matcher.Field) {
					self.incMissingFieldMap(matcher.Field, path)
					return fmt.Errorf(
						"Missing field mapping '%v' in %v", matcher.Field, logsource)
				}

				// Just report an unknown field but do not reject the
				// rule - this is just a warning that the rule is
				// using a field on this log source which is not known
				// to belong to this log source.
				pres := self.event_resolver.CheckFieldOnLogSource(
					matcher.Field, logsource)
				if !pres {
					self.incMissingFieldInLogSource(matcher.Field,
						logsource, path)
				}
			}
		}
	}

	details_any, pres := rule.AdditionalFields["details"]
	if pres {
		details, ok := details_any.(string)
		if ok {
			err = self.walk_details(details, path, logsource)
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

type LogSourceCollection struct {
	*ordereddict.Dict
}

func (self *LogSourceCollection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	dict := ordereddict.NewDict()

	mapping := make(map[string]Query)
	err := unmarshal(mapping)
	if err != nil {
		return err
	}

	err = unmarshal(dict)
	if err != nil {
		return err
	}

	for _, k := range dict.Keys() {
		v, pres := mapping[k]
		if pres {
			dict.Update(k, v)
		} else {
			dict.Delete(k)
		}
	}

	self.Dict = dict
	return nil
}

func NewLogSorceCollection() *LogSourceCollection {
	return &LogSourceCollection{
		Dict: ordereddict.NewDict(),
	}
}
