package main

type DefaultDetails struct {
	// A lambda that will be used to get the default description
	Query  string            `json:"Query"`
	Lookup map[string]string `json:"Lookup"`
}

type Query struct {
	Query   string   `json:"query"`
	Channel []string `json:"channel"`
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
	Preamble       string            `json:"Preamble"`
	FieldMappings  map[string]string `json:"FieldMappings"`
	DefaultDetails DefaultDetails    `json:"DefaultDetails"`
	Sources        map[string]Query  `json:"Sources"`

	// If this is set then we generate a reference URL for each rule.
	BaseReferenceURL string   `json:"BaseReferenceURL"`
	RuleDirectories  []string `json:"RuleDirectories"`
}
