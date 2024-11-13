package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func encode(in []byte) string {
	// Compress the string
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(in)
	gz.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (self *CompilerContext) getRules() []byte {
	return []byte(strings.Join(self.rules, "\n---\n"))
}

func (self *CompilerContext) GetArtifact() (string, error) {
	params := &ArtifactContent{
		Time:                       time.Now().UTC().Format(time.RFC3339),
		Base64CompressedRules:      encode(self.getRules()),
		Base64FieldMapping:         encode(MustMarshal(self.config_obj.FieldMappings)),
		Base64DefaultDetailsLookup: encode(MustMarshal(self.config_obj.DefaultDetails.Lookup)),
		Base64DefaultDetailsQuery:  strings.TrimSpace(self.config_obj.DefaultDetails.Query),
		LogSources:                 BuildLogSource(self.config_obj),
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

	return preamble_templ + export_templ + query_templ + postscript_temp, nil
}

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

	// Also include the redacted rules in the zip file
	fd, err = zip.Create("sigma_rules.yml")
	if err != nil {
		return err
	}
	fd.Write(self.getRules())

	fd, err = zip.Create("original_sigma_rules.yml")
	if err != nil {
		return err
	}
	fd.Write(self.original_rules.Bytes())

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
