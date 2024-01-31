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

func TestWalkFields_invalidLogsourceField(t *testing.T) {
	logsource := "windows/process_creation"
	correctName := "process_name"
	incorrectName := "processName"

	// Create a rule with a valid field, but invalid logsource field
	rule := CreateSimpleRule(correctName)

	context := CreateTestContext(correctName, incorrectName, logsource)

	err := context.walk_fields(&rule, "test", logsource)
	assert.Error(t, err)
}

func TestWalkFields_missingField(t *testing.T) {
	logsource := "windows/process_creation"
	correctName := "process_name"
	incorrectName := "processName"

	// Create a rule with a valid field, but invalid logsource field
	rule := CreateSimpleRule(correctName)

	context := CreateTestContext(incorrectName, incorrectName, logsource)

	err := context.walk_fields(&rule, "test", logsource)
	assert.Error(t, err)
}

func TestWalkFields_CorrectField(t *testing.T) {
	logsource := "windows/process_creation"
	correctName := "process_name"

	// Create a rule with a valid field, but invalid logsource field
	rule := CreateSimpleRule(correctName)

	context := CreateTestContext(correctName, correctName, logsource)

	err := context.walk_fields(&rule, "test", logsource)
	assert.NoError(t, err)
}
