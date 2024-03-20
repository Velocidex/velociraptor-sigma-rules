package main

import (
	"testing"

	"github.com/bradleyjkemp/sigma-go"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	description        string
	config             string
	rule               string
	logsource_err, err string

	// Additional tests
	context_check func(t *testing.T, context *CompilerContext)
}

var (
	test_cases = []testCase{
		{
			description: "Rule refers to fields and log sources.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
`,
			rule: `
logsource:
  product: windows
  service: test_service

detection:
  search:
    field_name|any:
       - XX
  condition: search
`,
		},
		{
			description: "Rule refers to fields not existing in field mappings.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
`,
			rule: `
logsource:
  product: windows
  service: test_service

detection:
  search:
    unspecified_field|any:
       - XX
  condition: search
`,
			err: "Missing field mapping 'unspecified_field'",
		},
		{
			description: "Rule refers to invalid log source.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
`,
			rule: `
logsource:
  product: windows
  service: another_service

detection:
  search:
    field_name|any:
       - XX
  condition: search
`,
			logsource_err: "Missing Source: '*/windows/another_service'",
		},
		{
			// Some rules define incorrect or very broad log sources
			// yet at the same time narrow down matches inside the
			// rule itself by checking against the event log
			// channel. If the config file specifies a channel for a
			// log source, we override the rule's log source selection
			// with the correct log source.
			description: "Rule refers to a log source but checks a channel that enables us to guess.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"
  Channel: "x=>x.Channel"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
   channel:
     - SomeChannel
`,
			// This rule specifies */windows/another_service as the
			// log source but this is inconsistant with the channel
			// check in the rule itself. We therefore override the
			// rule with */windows/test_service anyway.
			rule: `
logsource:
  product: windows
  service: another_service

detection:
  channel_check:
     Channel: SomeChannel

  search:
    field_name|any:
       - XX
  condition: search and channel_check
`,
		},
		{
			// The FieldMappings are global for all log sources. But
			// we need to be more specific than that in order to
			// detect broken rules. The config file can specify a list
			// of valid fields **for each log source**. We can use
			// this list to flag errors in the rule due to typos or
			// just invalid rules.
			description: "Checking fields per log source.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"
  invalid_field: "x=>x.invalid_field"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
   fields:
     - field_name
`,
			rule: `
logsource:
  product: windows
  service: test_service

detection:
  search:
    invalid_field: XXXX
    field_name|any:
       - XX
  condition: search
`,
			// The compile works but we need to flag some issues.
			context_check: func(t *testing.T, context *CompilerContext) {
				missing_fields := JsonDump(context.missing_fields_in_logsources)
				// This says that the test_service log source does not
				// have the "invalid_field" field which is used by
				// test.yml rule.
				assert.Equal(t, missing_fields,
					`{"*/windows/test_service":{"invalid_field":["/rules/test.yml"]}}`)
			},
		},
		{
			description: "Invalid field modifiers.",
			config: `
FieldMappings:
  field_name: "x=>x.field_name"

Sources:
 "*/windows/test_service":
   query: SELECT * FROM info()
`,
			rule: `
logsource:
  product: windows
  service: test_service

detection:
  search:
    field_name|somemodifier:
       - XX
  condition: search
`,
			err: "Invalid modifier somemodifier.",
		},
	}
)

func TestWalkFields(t *testing.T) {
	// Where the rule exists - for reporting errors etc
	file_path := "/rules/test.yml"

	for _, test_case := range test_cases {
		rule, err := sigma.ParseRule([]byte(test_case.rule))
		assert.NoError(t, err)

		context := NewCompilerContext()
		err = context.LoadConfigFromString(test_case.config)
		assert.NoError(t, err)

		// The log source is normalized by concatenating the product,
		// category and service
		logsource, err := context.normalize_logsource(&rule, file_path)
		if test_case.logsource_err == "" {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err, test_case.logsource_err)
		}

		// Walk the fields in the rule and check that they are present in
		// the field mappings.
		err = context.walk_fields(&rule, file_path, logsource)
		if test_case.err == "" {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err, test_case.err)
		}

		if test_case.context_check != nil {
			test_case.context_check(t, context)
		}
	}
}
