package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"text/template"
)

func encode(in []byte) string {
	// Compress the string
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(in)
	gz.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (self *CompilerContext) WriteArtifact(zip *zip.Writer) error {
	vql := BuildLogSource(self.config_obj)

	params := &ArtifactContent{
		Base64CompressedRules:      encode(self.rules.Bytes()),
		Base64FieldMapping:         encode(MustMarshal(self.config_obj.FieldMappings)),
		Base64DefaultDetailsLookup: encode(MustMarshal(self.config_obj.DefaultDetails.Lookup)),
		Base64DefaultDetailsQuery:  self.config_obj.DefaultDetails.Query,
	}

	templ, err := template.New("").Parse(self.config_obj.QueryTemplate)
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	err = templ.Execute(b, params)
	if err != nil {
		return err
	}

	vql += string(b.Bytes())

	fd, err := zip.Create("artifact.yaml")
	if err != nil {
		return err
	}
	fd.Write([]byte(self.config_obj.Preamble + indent(vql, 4)))

	// Also include the redacted rules in the zip file
	fd, err = zip.Create("sigma_rules.yml")
	if err != nil {
		return err
	}
	fd.Write(self.rules.Bytes())

	fd, err = zip.Create("original_sigma_rules.yml")
	if err != nil {
		return err
	}
	fd.Write(self.original_rules.Bytes())

	fd, err = zip.Create("field_mapping.json")
	if err != nil {
		return err
	}
	fd.Write(MustMarshal(self.config_obj.FieldMappings))

	fd, err = zip.Create("default_details.json")
	if err != nil {
		return err
	}
	fd.Write(MustMarshal(self.config_obj.DefaultDetails.Lookup))

	return nil
}
