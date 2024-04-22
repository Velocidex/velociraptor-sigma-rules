package main

import (
	"fmt"
	"strings"
)

func BuildLogSource(config_obj *Config) string {
	var sources []string

	for k, v := range config_obj.Sources {
		query := strings.TrimSpace(v.Query)
		if len(query) > 0 {
			sources = append(sources, fmt.Sprintf(
				"  `%v`={\n%v}", k, indent(query, 5)))
		}
	}

	return "LET LogSources <= sigma_log_sources(\n" +
		strings.Join(sources, ",\n") +
		"\n)\n"
}

type ArtifactContent struct {
	Time                       string
	Base64CompressedRules      string
	Base64FieldMapping         string
	Base64DefaultDetailsLookup string
	Base64DefaultDetailsQuery  string
}
