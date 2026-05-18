package main

import (
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

type compileTestCase struct {
	config   string
	rule     string
	compiled string
}

var compileTestCases = []compileTestCase{
	{
		config: `
FieldMappings:
  Channel: "x=>x.Channel"

Sources:
 "*/windows/first_service":
   query: SELECT * ...
   channel:
     - SomeChannel
`,
		// Rule specifies a log source but makes a Channel check to
		// SomeChannel. The logsource specified is wrong and the rule
		// is rewritten to refer to the real log source.

		// This allows us to normalize rules which refer to a generic
		// log source specifically but then post filter all the
		// results on the channel, when there is a more specific log
		// source available that could be used already.
		rule: `
title: Update log source based on Channel check
logsource:
  product: windows
  service: second_service

detection:
  channel_check:
     Channel: SomeChannel
  condition: channel_check
`,
	},
	{
		config: `
FieldMappings:
  Channel: "x=>x.Channel"
  SomeField: "x=>x.SomeField"

Sources:
 "*/windows/first_service":
   query: SELECT * ...
   channel:
     - SomeChannel
`,
		// Sometimes a rule will not specify a condition (Is this even
		// valid Sigma?). In that case we should AND all the
		// detections
		rule: `
title: No condition present in rule
logsource:
  product: windows
  service: first_service

detection:
  channel_check:
     Channel: SomeChannel
  selection:
     SomeField: XXXX

`,
	},
}

func TestCompilation(t *testing.T) {
	golden := []string{}

	for _, test_case := range compileTestCases {
		context := NewCompilerContext()
		err := context.LoadConfigFromString(test_case.config)
		assert.NoError(t, err)

		err = context.CompileRule(test_case.rule, "/path/to/rule.yml")
		assert.NoError(t, err)

		golden = append(golden, string(context.getRules()))
	}

	g := goldie.New(
		t,
		goldie.WithFixtureDir("fixtures"),
		goldie.WithDiffEngine(goldie.ClassicDiff),
	)

	g.Assert(t, "TestCompilation", []byte(strings.Join(golden, "\n---\n")))
}

func TestCorrelationCompilation(t *testing.T) {
	config := `
FieldMappings:
  CommandLine: "x=>x.CommandLine"

Sources:
 "process_creation/linux/*":
   query: SELECT * FROM info()
`

	source_rule_one := `
title: First Source Rule
name: source_rule_one
id: 11111111-1111-1111-1111-111111111111
logsource:
  product: linux
  category: process_creation
detection:
  selection:
    CommandLine: foo
  condition: selection
`

	source_rule_two := `
title: Second Source Rule
name: source_rule_two
id: 22222222-2222-2222-2222-222222222222
logsource:
  product: linux
  category: process_creation
detection:
  selection:
    CommandLine: bar
  condition: selection
`

	correlation_rule := `
title: Correlation Of Two Source Rules
name: correlation_example
id: 33333333-3333-3333-3333-333333333333
correlation:
    type: temporal
    rules:
        - source_rule_one
        - source_rule_two
    group-by:
        - Computer
    timespan: 5m
level: high
`

	context := NewCompilerContext()
	err := context.LoadConfigFromString(config)
	assert.NoError(t, err)

	for _, rule := range []string{source_rule_one, source_rule_two, correlation_rule} {
		err = context.CompileRule(rule, "/path/to/rule.yml")
		assert.NoError(t, err)
	}

	rendered := string(context.getRules())

	// Spot-check patch invariants before goldie: correlation block and
	// source-rule references must survive compilation, else the patch regressed.
	assert.Contains(t, rendered, "correlation:")
	assert.Contains(t, rendered, "source_rule_one")
	assert.Contains(t, rendered, "source_rule_two")
	assert.Contains(t, rendered, "correlation_example")
	assert.Contains(t, rendered, "11111111-1111-1111-1111-111111111111")

	g := goldie.New(
		t,
		goldie.WithFixtureDir("fixtures"),
		goldie.WithDiffEngine(goldie.ClassicDiff),
	)

	g.Assert(t, "TestCorrelationCompilation", []byte(rendered))
}
