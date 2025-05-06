package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func encode(in []byte) string {
	// Compress the string
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(in)
	gz.Close()
	res := base64.StdEncoding.EncodeToString(b.Bytes())
	return res
}

func (self *CompilerContext) getRules() []byte {
	return []byte(strings.Join(self.rules, "\n---\n"))
}

func (self *CompilerContext) GetArtifact() (string, error) {
	if self.completed_artifact != "" {
		return self.completed_artifact, nil
	}

	params := &ArtifactContent{
		Time:                       time.Now().UTC().Format(time.RFC3339),
		Base64CompressedRules:      encode(self.getRules()),
		Base64FieldMapping:         encode(MustMarshal(self.config_obj.FieldMappings)),
		Base64DefaultDetailsLookup: encode(MustMarshal(self.config_obj.DefaultDetails.Lookup)),
		Base64DefaultDetailsQuery:  strings.TrimSpace(self.config_obj.DefaultDetails.Query),
		LogSources:                 BuildLogSource(self.config_obj),
	}

	for _, imp := range self.imported_configs {
		params.ImportedLogSources = append(params.ImportedLogSources,
			BuildLogSource(imp)...)
	}

	// Allow the artifact to export functions to other artifacts.
	export_templ, err := calculateTemplate(self.config_obj.ExportTemplate, params)
	if err != nil {
		return "", err
	}

	preamble_templ, err := calculateTemplate(self.config_obj.Preamble, params)
	if err != nil {
		return "", err
	}

	query_templ, err := calculateTemplate(self.config_obj.QueryTemplate, params)
	if err != nil {
		return "", err
	}

	postscript_temp, err := calculateTemplate(self.config_obj.Postscript, params)
	if err != nil {
		return "", err
	}

	self.completed_artifact = preamble_templ + export_templ + query_templ + postscript_temp

	return self.completed_artifact, nil
}

/*
Build the artifact zip file.

The artifact zip file contains all the sigma rules as well as the
raw artifact, and a script which can be used to rebuild the artifact
based on a subset of the rules.

Users can unpack the zip file, remove noisy rules and rebuild the
artifact easily.
*/
func (self *CompilerContext) WriteArtifact(zip *zip.Writer) error {
	artifact_yaml, err := self.GetArtifact()
	if err != nil {
		return err
	}
	artifact_name := self.config_obj.Name
	if artifact_name == "" {
		artifact_name = "artifact"
	}

	fd, err := zip.Create(artifact_name + ".yaml")
	if err != nil {
		return err
	}
	fd.Write([]byte(artifact_yaml))

	// Add the rules into the zip file.
	err = self.WriteRules(zip)
	if err != nil {
		return err
	}

	fd, err = zip.Create("field_mapping.json")
	if err != nil {
		return err
	}
	fd.Write(MustMarshal(self.config_obj.field_mappings))

	fd, err = zip.Create("default_details.json")
	if err != nil {
		return err
	}
	fd.Write(MustMarshal(self.config_obj.DefaultDetails.Lookup))

	fd, err = zip.Create("rejected.json")
	if err != nil {
		return err
	}
	fd.Write(MustMarshalIndent(self.GetRejected()))

	for _, inc := range self.config_obj.IncludeArtifacts {
		fd, err := os.Open(inc)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		fd.Close()

		out_fd, err := zip.Create(filepath.Base(inc))
		if err != nil {
			return err
		}

		out_fd.Write(data)
	}

	return nil
}

func (self *CompilerContext) WriteRules(zip *zip.Writer) error {

	for rule_path, original_rule := range self.original_rules_by_path {
		target := path.Join("original_rules", rule_path)

		fd, err := zip.Create(target)
		if err != nil {
			return err
		}
		fd.Write([]byte(original_rule))

		rule, pres := self.rules_by_path[rule_path]
		if pres {
			target := path.Join("rules", rule_path)
			fd, err := zip.Create(target)
			if err != nil {
				return err
			}

			buf := &bytes.Buffer{}
			yamlEncoder := yaml.NewEncoder(buf)
			yamlEncoder.SetIndent(2)
			err = yamlEncoder.Encode(rule)
			if err != nil {
				return err
			}

			fd.Write(buf.Bytes())
		}
	}

	fd, err := zip.Create("compiler.py")
	if err != nil {
		return err
	}
	fd.Write([]byte(CompierPython))

	return nil
}
