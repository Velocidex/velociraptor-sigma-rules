package main

type DefaultDetails struct {
	// A lambda that will be used to get the default description
	Query  string            `json:"Query"`
	Lookup map[string]string `json:"Lookup"`
}

type Query struct {
	Query       string   `json:"query"`
	Channel     []string `json:"channel"`
	Fields      []string `json:"fields"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
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
	Name           string            `json:"Name"`
	Preamble       string            `json:"Preamble"`
	FieldMappings  map[string]string `json:"FieldMappings"`
	DefaultDetails DefaultDetails    `json:"DefaultDetails"`
	Sources        map[string]Query  `json:"Sources"`
	ExportTemplate string            `json:"ExportTemplate"`
	QueryTemplate  string            `json:"QueryTemplate"`
	Postscript     string            `json:"Postscript"`

	// If this is set then we generate a reference URL for each rule.
	BaseReferenceURL string   `json:"BaseReferenceURL"`
	RuleDirectories  []string `json:"RuleDirectories"`

	// Many rules are broken and have bad field mappings or log
	// sources. The following list suppresses these warnings (but the
	// rules are still rejected)
	BadFieldMappings []string `json:"BadFieldMappings"`
	BadSources       []string `json:"BadSources"`

	EventResolver string `json:"EventResolver"`

	// Include these artifacts into the zip bundle as well. There are
	// relative paths to the included files. These are usually used to
	// include dependent artifacts.
	IncludeArtifacts []string `json:"IncludeArtifacts"`

	// Read these configs also. Many attributes are merged with
	// included configs (for example FieldMappings, and Sources). This
	// allows to build derived artifacts based on other artifacts.
	ImportConfigs []string `json:"ImportConfigs"`

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
