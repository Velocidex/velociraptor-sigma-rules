package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func (self *CompilerContext) GetDocs() (string, error) {
	params := &ArtifactContent{
		Time:       time.Now().UTC().Format(time.RFC3339),
		LogSources: BuildLogSource(self.config_obj),
		Config:     self.config_obj,
	}

	for k, v := range self.config_obj.FieldMappings {
		params.FieldMappings = append(params.FieldMappings, FieldMapping{
			Name: k, Mapping: v,
		})
	}

	sort.Slice(params.FieldMappings, func(i, j int) bool {
		return params.FieldMappings[i].Name < params.FieldMappings[j].Name
	})

	templ, err := calculateTemplate(self.config_obj.DocTemplate, params)
	if err != nil {
		return "", err
	}

	return templ, nil
}

type RuleIndex struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Author      string   `json:"author,omitempty"`
	Link        string   `json:"link,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (self *CompilerContext) WriteRuleDir(base string) error {
	var index []RuleIndex

	for path, original_rule := range self.original_rules_by_path {
		target := filepath.Join(base, path)
		err := os.MkdirAll(filepath.Dir(target), 0700)
		if err != nil {
			return err
		}

		fd, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		fd.Write([]byte(original_rule))
		fd.Close()
		fmt.Printf("Wrote rule %v\n", target)

		rule, pres := self.rules_by_path[path]
		if pres {
			index = append(index, RuleIndex{
				Title:       rule.Title,
				Description: rule.Description,
				Author:      rule.Author,
				Tags:        rule.Tags,
				Link:        path,
			})
		}
	}

	serialized, err := json.Marshal(index)
	if err != nil {
		return err
	}

	fd, err := os.OpenFile(filepath.Join(base, "index.json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	fd.Write(serialized)
	return fd.Close()
}
