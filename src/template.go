package main

import (
	"bytes"
	"sort"
	"strings"
	"text/template"
)

func BuildLogSource(config_obj *Config) []Query {
	var sources []Query

	// If this artifact does not add any new log sources, we do not
	// need to emit any log sources, as we will use the one from the
	// export section of the base artifact.
	if len(config_obj.Sources) == 0 {
		return nil
	}

	if config_obj.sources == nil {
		config_obj.mergeConfig(config_obj)
	}

	for k, v := range config_obj.sources {
		query := strings.TrimSpace(v.Query)
		if len(query) > 0 {
			sources = append(sources, Query{
				Query:       query,
				Name:        k,
				Description: v.Description})
		}
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].Name < sources[j].Name
	})

	return sources
}

type ArtifactContent struct {
	Time                       string
	Base64CompressedRules      string
	Base64FieldMapping         string
	Base64DefaultDetailsLookup string
	Base64DefaultDetailsQuery  string
	LogSources                 []Query
	ImportedLogSources         []Query
}

func indentTemplate(args ...interface{}) interface{} {
	if len(args) != 2 {
		return ""
	}

	template, ok := args[0].(string)
	if !ok {
		return ""
	}

	indent_size, ok := args[1].(int)
	if !ok {
		return template
	}

	return indent(template, indent_size)
}

func calculateTemplate(template_str string, params *ArtifactContent) (string, error) {
	templ, err := template.New("").Funcs(
		template.FuncMap{
			"Indent": indentTemplate,
		}).Parse(template_str)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = templ.Execute(b, params)
	if err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}
