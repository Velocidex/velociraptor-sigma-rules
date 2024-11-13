package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := map[string]struct {
		configFile            string
		configLoadError       string
		doCompileError        string
		expectedMissingFields int
		expectedInvalidFields int
	}{
		"CuratedRulesTest": {
			configFile:            "../config/windows_base.yaml",
			configLoadError:       "",
			doCompileError:        "",
			expectedMissingFields: 0,
			expectedInvalidFields: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			context := NewCompilerContext()
			level_regex, err := regexp.Compile(".")
			if err != nil {
				t.Fatalf("Level Regex invalid: %v", err)
			}
			context.level_regex = level_regex

			err = context.LoadConfig(tc.configFile)
			if err == nil {
				if tc.configLoadError != "" {
					t.Fatalf("Expected error %v, got none", tc.configLoadError)
				}
			} else {
				if tc.configLoadError == "" {
					t.Fatalf("Expected no error, got %v", err)
				} else {
					assert.Contains(t, err.Error(), tc.configLoadError)
				}
			}

			// Fix up the rule path to make it valid as Tests execute out of project path.
			for i, path := range context.config_obj.RuleDirectories {
				context.config_obj.RuleDirectories[i] = "." + path
			}

			err = context.CompileDirs()
			if err == nil {
				if tc.doCompileError != "" {
					t.Errorf("Expected error %v, got none", tc.doCompileError)
				}
			} else {
				if tc.doCompileError == "" {
					t.Errorf("Expected no error, got %v", err)
				} else {
					assert.Contains(t, err.Error(), tc.doCompileError)
				}
			}

			context.Stats()

			if len(context.missing_fields) != tc.expectedMissingFields {
				t.Errorf("Expected %v missing fields, got %v", tc.expectedMissingFields, len(context.missing_fields))
			}

			if len(context.invalid_fields) != tc.expectedInvalidFields {
				t.Errorf("Expected %v invalid fields, got %v", tc.expectedInvalidFields, len(context.invalid_fields))
			}
		})
	}
}
