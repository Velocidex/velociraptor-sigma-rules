package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	common_fields = []string{"Channel", "EventID", ""}
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

	log_def, pres := self.config_obj.Sources[logsource]
	if !pres {
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
	_, pres := self.config_obj.FieldMappings[field]
	if pres {
		return true
	}

	// Is the field mapping known by any of our log sources?
	_, pres = self.all_fields[field]
	if pres {
		// Add an automatic log source
		self.config_obj.FieldMappings[field] = fmt.Sprintf("x=>x.EventData.%s", field)
		return true
	}

	return false
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
