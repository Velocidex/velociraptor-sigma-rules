package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Velocidex/sigma-go"
)

func BuildLogSource(config_obj *Config) []Query {
	var sources []Query

	// If this artifact does not add any new log sources, we do not
	// need to emit any log sources, as we will use the one from the
	// export section of the base artifact.
	if config_obj.Sources.Len() == 0 {
		return nil
	}

	if config_obj.sources == nil || config_obj.sources.Len() == 0 {
		config_obj.mergeConfig(config_obj)
	}

	for _, k := range config_obj.sources.Keys() {
		query_any, pres := config_obj.sources.Get(k)
		if !pres {
			continue
		}

		query, ok := query_any.(Query)
		if !ok {
			continue
		}

		query_str := strings.TrimSpace(query.Query)
		if query.Summary == "" {
			lines := strings.Split(query.Description, "\n")
			for _, l := range lines {
				if len(l) > 0 {
					query.Summary = l
					break
				}
			}
		}

		if len(query_str) > 0 {
			q := Query{
				Query:       query_str,
				Name:        k,
				Description: query.Description,
				Summary:     query.Summary,
				LogSource:   &sigma.Logsource{},
				Samples:     query.Samples,
			}
			updateRuleLogSources(k, q.LogSource)

			sources = append(sources, q)
		}
	}

	return sources
}

type FieldMapping struct {
	Name    string
	Mapping string
}

type ArtifactContent struct {
	Time                       string
	Base64CompressedRules      string
	Base64FieldMapping         string
	Base64DefaultDetailsLookup string
	Base64DefaultDetailsQuery  string
	LogSources                 []Query
	ImportedLogSources         []Query
	FieldMappings              []FieldMapping
	Config                     *Config
}

func readFile(args ...interface{}) interface{} {
	result := ""

	for _, arg := range args {
		path, ok := arg.(string)
		if !ok {
			continue
		}

		fd, err := os.Open(path)
		if err != nil {
			continue
		}
		defer fd.Close()

		data, err := ioutil.ReadAll(fd)
		if err != nil {
			continue
		}

		result += string(data)
	}

	return result
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
			"Indent":   indentTemplate,
			"ReadFile": readFile,
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
