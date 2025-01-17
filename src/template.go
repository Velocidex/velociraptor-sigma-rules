package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Velocidex/sigma-go"
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
			q := Query{
				Query:       query,
				Name:        k,
				Description: v.Description,
				LogSource:   &sigma.Logsource{},
				Samples:     v.Samples,
			}
			updateRuleLogSources(k, q.LogSource)

			sources = append(sources, q)
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
