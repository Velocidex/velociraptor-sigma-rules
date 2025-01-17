package main

import "github.com/Velocidex/sigma-go"

type DefaultDetails struct {
	// A lambda that will be used to get the default description
	Query  string            `json:"Query,omitempty"`
	Lookup map[string]string `json:"Lookup,omitempty"`
}

type Sample struct {
	Name string `json:"name,omitempty"`
	Json string `json:"json,omitempty"`
}

type Query struct {
	Query       string           `json:"query,omitempty"`
	Channel     []string         `json:"channel,omitempty"`
	Fields      []string         `json:"fields,omitempty"`
	Description string           `json:"description,omitempty"`
	Name        string           `json:"name,omitempty"`
	LogSource   *sigma.Logsource `json:"log_source,omitempty"`

	// A set of JSON files with sample events for this log source
	Samples []Sample `json:"samples,omitempty"`
}

// Specify source transformations.

// Sometimes sigma rules refer to the same source differently, either
// as an alias or as a filtered version of the same source. In order
// to ensure we only read the same source once, we need to transform
// the rule to read from the same source.
type SourceRemapping struct {
	SubstituteSource string `json:"SubstituteSource"`
}

type Config struct {
	Name           string            `json:"Name,omitempty"`
	Description    string            `json:"Description,omitempty"`
	Preamble       string            `json:"Preamble,omitempty"`
	FieldMappings  map[string]string `json:"FieldMappings,omitempty"`
	DefaultDetails DefaultDetails    `json:"DefaultDetails,omitempty"`
	Sources        map[string]Query  `json:"Sources,omitempty"`
	ExportTemplate string            `json:"ExportTemplate,omitempty"`
	DocTemplate    string            `json:"DocTemplate,omitempty"`
	QueryTemplate  string            `json:"QueryTemplate,omitempty"`
	Postscript     string            `json:"Postscript,omitempty"`

	// If this is set then we generate a reference URL for each rule.
	BaseReferenceURL string   `json:"BaseReferenceURL,omitempty"`
	RuleDirectories  []string `json:"RuleDirectories,omitempty"`

	// Many rules are broken and have bad field mappings or log
	// sources. The following list suppresses these warnings (but the
	// rules are still rejected)
	BadFieldMappings []string `json:"BadFieldMappings,omitempty"`
	BadSources       []string `json:"BadSources,omitempty"`

	EventResolver string `json:"EventResolver,omitempty"`

	// Include these artifacts into the zip bundle as well. There are
	// relative paths to the included files. These are usually used to
	// include dependent artifacts.
	IncludeArtifacts []string `json:"IncludeArtifacts,omitempty"`

	// Read these configs also. Many attributes are merged with
	// included configs (for example FieldMappings, and Sources). This
	// allows to build derived artifacts based on other artifacts.
	ImportConfigs []string `json:"ImportConfigs,omitempty"`

	// Merged results from imported configs
	sources        map[string]Query
	field_mappings map[string]string
}

// Merge fields from the config into this state object.
func (self *Config) mergeConfig(config_obj *Config) {
	if self.sources == nil {
		self.sources = make(map[string]Query)
	}

	if config_obj.Sources != nil {
		for k, v := range config_obj.Sources {
			self.sources[k] = v
		}
	}

	if self.field_mappings == nil {
		self.field_mappings = make(map[string]string)
	}

	if config_obj.FieldMappings != nil {
		for k, v := range config_obj.FieldMappings {
			self.field_mappings[k] = v
		}
	}
}
