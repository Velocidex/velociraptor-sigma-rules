package main

import (
	"testing"

	"github.com/bradleyjkemp/sigma-go"
	"github.com/stretchr/testify/assert"
)

func CreateSimpleRule(name string) sigma.Rule {
	field := &sigma.FieldMatcher{
		Field: name,
	}

	detection := sigma.Detection{
		Searches: map[string]sigma.Search{
			"search": sigma.Search{
				EventMatchers: []sigma.EventMatcher{
					sigma.EventMatcher{
						*field,
					},
				},
			},
		},
	}

	return sigma.Rule{
		Title:     "Test Rule",
		Detection: detection,
	}
}

func CreateTestContext(fieldName string, logsourceFieldName string, logsource string) *CompilerContext {
	context := NewCompilerContext()
	config := &Config{
		FieldMappings: map[string]string{
			fieldName: "test",
		},
		Sources: map[string]Query{
			logsource: Query{
				Fields: []string{logsourceFieldName},
			},
		},
	}
	context.config_obj = config
	return context
}

func TestWalkFields(t *testing.T) {
	logsource := "windows/process_creation"
	correctName := "process_name"
	incorrectName := "processName"

	tests := map[string]struct {
		RuleField      string
		ConfigField    string
		LogsourceField string
		Logsource      string
		Want           string
	}{
		"Correct Field":             {RuleField: correctName, ConfigField: correctName, LogsourceField: correctName, Logsource: logsource, Want: ""},
		"Invalid Field":             {RuleField: correctName, ConfigField: incorrectName, LogsourceField: incorrectName, Logsource: logsource, Want: "Missing field"},
		"Incorrect Logsource Field": {RuleField: correctName, ConfigField: correctName, LogsourceField: incorrectName, Logsource: logsource, Want: "Invalid field"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rule := CreateSimpleRule(tc.RuleField)

			context := CreateTestContext(tc.ConfigField, tc.LogsourceField, tc.Logsource)

			err := context.walk_fields(&rule, "test", logsource)
			if tc.Want == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tc.Want)
			}
		})
	}
}
