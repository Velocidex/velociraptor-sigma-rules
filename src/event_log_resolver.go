package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	common_fields = []string{"Channel", "EventID", ""}

	// The event object passed to a field-mapping lambda exposes these
	// top-level members. A mapping body must reach them through the lambda
	// argument (e.g. "x=>x.System.UserID"); referencing one as a bare free
	// symbol (e.g. "x=>System.UserID") compiles but throws at evaluation
	// time with "Symbol System not found". This regexp matches a bare
	// reference - a root that starts a fresh path expression rather than
	// following "x.".
	bare_event_root_regexp = regexp.MustCompile(
		`(?:^|[(=,>\s])(System|EventData|UserData|Message)\.`)
)

type EventSchema struct {
	Id      string   `json:"Id"`
	Channel string   `json:"Channel"`
	Message string   `json:"Message"`
	Fields  []string `json:"Fields"`
}

type EventResolver struct {
	schema []EventSchema

	config_obj *Config

	// All the fields in each channel: key: lower(channel), value list of fields
	field_by_channel map[string][]string

	// All known fields:
	all_fields map[string]bool
}

func (self *EventResolver) CheckFieldOnLogSource(field, logsource string) bool {
	if slices.Contains(common_fields, field) {
		return true
	}

	log_def_any, pres := self.config_obj.sources.Get(logsource)
	if !pres {
		return false
	}

	log_def, ok := log_def_any.(Query)
	if !ok {
		return false
	}

	if slices.Contains(log_def.Fields, field) {
		return true
	}

	check_channel := func(channel string) bool {
		channel_key := strings.ToLower(channel)
		fields, pres := self.field_by_channel[channel_key]
		if !pres {
			return false
		}

		for _, f := range fields {
			if f == field {
				return true
			}
		}

		return false
	}

	for _, channel := range log_def.Channel {
		if check_channel(channel) {
			return true
		}
	}

	return false
}

func (self *EventResolver) CheckFieldMapping(field string) bool {
	// Is the fields mapping defined in the config file?
	_, pres := self.config_obj.field_mappings[field]
	if pres {
		return true
	}

	// If the field has "." it might be a compound field
	if strings.Contains(field, ".") {
		return true
	}

	// Add an automatic log source
	fmt.Printf("Error: Need to add the following field mapping to the base artifact:\n %v: \"x=>x.EventData.%s\"\n",
		field, field)
	return false
}

// CheckFieldMappingLambda returns the offending event-root field name if
// the field-mapping lambda references it as a bare free symbol instead of
// reaching it through the lambda argument (x.<root>), or "" if the lambda
// is well formed. This catches the typo class "x=>System.UserID" (which
// should be "x=>x.System.UserID") that compiles but fails at evaluation
// time. Function-call lambdas (e.g. "x=>timestamp(epoch=now())") and
// roots reached through the argument or nested in a function argument
// (e.g. "x=>process_tracker_get(id=x.System.ProcessID).Data.Exe") are not
// flagged.
func CheckFieldMappingLambda(lambda string) string {
	// Only inspect the lambda body, after the "=>".
	body := lambda
	if _, after, found := strings.Cut(lambda, "=>"); found {
		body = after
	}

	matches := bare_event_root_regexp.FindStringSubmatch(body)
	if matches == nil {
		return ""
	}
	return matches[1]
}

// CheckFieldMappings validates every field-mapping lambda in the loaded
// config and reports any that reference an event-root field as a bare free
// symbol. Following the same non-fatal reporting style as CheckFieldMapping,
// it prints an Error line (surfaced in compile/CI output) rather than
// aborting the build.
func (self *CompilerContext) CheckFieldMappings() {
	for name, lambda := range self.config_obj.field_mappings {
		root := CheckFieldMappingLambda(lambda)
		if root == "" {
			continue
		}

		fmt.Printf("Error: Field mapping %q (%v) references %q as a free symbol; "+
			"did you mean x.%v? VQL fails with \"Symbol %v not found\" at evaluation time.\n",
			name, lambda, root, root, root)
	}
}

func (self *EventResolver) Load(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	self.schema = nil
	self.field_by_channel = make(map[string][]string)
	self.all_fields = make(map[string]bool)

	reader := bufio.NewReader(fd)
	for {
		row_data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		// We have reached the end.
		if len(row_data) == 0 {
			break
		}

		if len(row_data) < 2 {
			continue
		}

		item := EventSchema{}
		err = json.Unmarshal(row_data, &item)
		if err != nil {
			continue
		}
		self.schema = append(self.schema, item)
		channel_key := strings.ToLower(item.Channel)
		existing, _ := self.field_by_channel[channel_key]
		existing = append(existing, item.Fields...)
		self.field_by_channel[channel_key] = existing

		for _, f := range item.Fields {
			self.all_fields[f] = true
		}
	}

	fmt.Printf("Loaded Event Resolver with %v definitions\n", len(self.schema))

	return nil
}
