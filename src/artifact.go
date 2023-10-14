package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
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

	vql += fmt.Sprintf(`
LET Rules <= gunzip(string=base64decode(string=%q))
LET FieldMapping <= parse_json(data=gunzip(string=base64decode(string=%q)))
LET DefaultDetails <= parse_json(data=gunzip(string=base64decode(string=%q)))
LET X = scope()

SELECT timestamp(epoch=System.TimeCreated.SystemTime) AS Timestamp,
       System.Computer AS Computer,
       System.Channel AS Channel,
       System.EventID.Value AS EID,
       _Rule.Level AS Level,
       _Rule.Title AS Title,
       System.EventRecordID AS RecordID,
       Details,
       dict(System=System, EventData=X.EventData || X.UserData, Message=X.Message) AS _Event
 FROM sigma(
   rules=split(string= Rules, sep_string="\n---\n"),
   log_sources= LogSources, debug=Debug,
   default_details=%q,
   rule_filter="x=>x.Level =~ RuleLevelRegex AND x.Status =~ RuleStatusRegex AND x.Title =~ RuleTitleFilter",
   field_mapping= FieldMapping)
`,
		encode(self.rules.Bytes()),
		encode(MustMarshal(self.config_obj.FieldMappings)),
		encode(MustMarshal(self.config_obj.DefaultDetails.Lookup)),
		self.config_obj.DefaultDetails.Query,
	)

	fd, err := zip.Create("artifact.yaml")
	if err != nil {
		return err
	}
	fd.Write([]byte(self.config_obj.Preamble + indent(vql, 4)))

	fd, err = zip.Create("sigma_rules.yml")
	if err != nil {
		return err
	}
	fd.Write(self.rules.Bytes())

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
