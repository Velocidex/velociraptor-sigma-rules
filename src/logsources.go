package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Velocidex/yaml/v2"
	"github.com/bradleyjkemp/sigma-go"
)

type CompilerContext struct {
	logsources     map[string]int
	fields         map[string]int
	missing_fields map[string][]string

	errored_rules map[string][]string
	config_obj    *Config

	rules       bytes.Buffer
	level_regex *regexp.Regexp

	vql bytes.Buffer
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		logsources:    make(map[string]int),
		errored_rules: make(map[string][]string),

		fields:         make(map[string]int),
		missing_fields: make(map[string][]string),
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

	config_obj := &Config{}
	err = yaml.Unmarshal(data, config_obj)
	if err != nil {
		return err
	}
	self.config_obj = config_obj

	return nil
}

func (self *CompilerContext) Resolve(source_spec string) bool {
	_, pres := self.config_obj.Sources[source_spec]
	if !pres {
		return false
	}
	count, _ := self.logsources[source_spec]
	count++

	self.logsources[source_spec] = count
	return true
}

func (self *CompilerContext) Stats() {
	sources := []string{}
	for k := range self.logsources {
		sources = append(sources, k)
	}

	sort.Strings(sources)
	fmt.Printf("The following log sources will be used:\n")
	for _, v := range sources {
		fmt.Printf("  %v\n", v)
	}

	if len(self.errored_rules) > 0 {
		fmt.Printf("\nErrored Rules without a valid log soure:\n")
		for k, v := range self.errored_rules {
			fmt.Printf("  %v:\n", k)
			for _, i := range v {
				fmt.Printf("     %v\n", i)
			}
		}
	}

	if len(self.missing_fields) > 0 {
		fmt.Printf("\nMissing field mappings:\n")
		for k, v := range self.missing_fields {
			fmt.Printf("  %v:\n", k)
			for _, i := range v {
				fmt.Printf("     %v\n", i)
			}
		}
	}
}

func (self *CompilerContext) getSourceFromChannel(channel string) string {
	for source, query := range self.config_obj.Sources {
		for _, c := range query.Channel {
			if c == channel {
				return source
			}
		}
	}
	return ""
}

func (self *CompilerContext) guessLogSource(rule *sigma.Rule, category, product, service string) (string, string) {
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
								parts := strings.Split(source, "/")
								if len(parts) == 3 {
									return source, parts[2]
								}
							}
						}
					}
				}
			}
		}
	}

	return fmt.Sprintf("%v/%v/%v", category, product, service), service
}

func (self *CompilerContext) walk_fields(rule *sigma.Rule, path string) error {
	for _, search := range rule.Detection.Searches {
		for _, event_matcher := range search.EventMatchers {
			for _, matcher := range event_matcher {
				// Check if there is a field mapping
				_, pres := self.config_obj.FieldMappings[matcher.Field]
				if !pres {
					missing, _ := self.missing_fields[matcher.Field]
					missing = append(missing, path)
					self.missing_fields[matcher.Field] = missing
					return errors.New("Missing field")
				}
			}
		}
	}
	return nil
}

func (self *CompilerContext) normalize_logsource(rule *sigma.Rule, path string) string {
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
	if !self.Resolve(source_spec) {
		// Try to guess the source spec
		guessed_source_spec, guessed_service := self.guessLogSource(rule, category, product, service)
		if !self.Resolve(guessed_source_spec) {
			fmt.Printf("**** Log Source '%v' not found!\n", guessed_source_spec)
			failed, _ := self.errored_rules[guessed_source_spec]
			failed = append(failed, path)
			self.errored_rules[guessed_source_spec] = failed
		} else {
			fmt.Printf(
				"** Substitute guess source %v for %v with rule %v\n",
				guessed_source_spec, source_spec, path)
			rule.Logsource.Service = guessed_service
			service = guessed_service
		}
	}

	return source_spec
}
