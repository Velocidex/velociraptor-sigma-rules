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
