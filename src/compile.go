package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/bradleyjkemp/sigma-go"
	"gopkg.in/yaml.v3"
)

var (
	app = kingpin.New("velosigma",
		"A tool for manipulating sigma files.")

	compile_cmd  = app.Command("compile", "Compile all the rules into one rule.")
	compile_dirs = compile_cmd.Arg("directory", "Directory to compile").Required().Strings()
	output_file  = compile_cmd.Flag("output", "File to write").Required().String()
	base_url     = compile_cmd.Flag("base_url", "URL to rule").
			Default("https://github.com/Yamato-Security/hayabusa-rules/tree/main/").
			String()

	command_handlers []CommandHandler
)

type CommandHandler func(command string) bool

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	for _, handler := range command_handlers {
		if handler(command) {
			break
		}
	}
}

func doCompile() {
	out_fd, err := os.OpenFile(*output_file,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	kingpin.FatalIfError(err, "Creating output")

	for _, compile_dir := range *compile_dirs {
		err := filepath.WalkDir(compile_dir,
			func(path string, d fs.DirEntry, err error) error {
				if !strings.HasSuffix(path, ".yml") {
					return nil
				}

				fd, err := os.Open(path)
				if err != nil {
					return err
				}

				data, err := ioutil.ReadAll(fd)
				if err != nil {
					return err
				}

				rule, err := sigma.ParseRule(data)
				if err != nil {
					return nil
				}

				additional_fields := make(map[string]interface{})
				details := rule.AdditionalFields["details"]
				if details != nil {
					additional_fields["details"] = details
				}

				new_rule := sigma.Rule{
					Title:            rule.Title,
					Author:           rule.Author,
					Level:            rule.Level,
					Logsource:        rule.Logsource,
					Detection:        rule.Detection,
					AdditionalFields: additional_fields,
					References:       []string{*base_url + path},
				}

				yamlEncoder := yaml.NewEncoder(out_fd)
				yamlEncoder.SetIndent(2)
				err = yamlEncoder.Encode(new_rule)
				if err != nil {
					return nil
				}

				fmt.Printf("Processing %v\n", path)
				out_fd.Write([]byte("\n---\n"))

				return err
			})
		kingpin.FatalIfError(err, "Listing directory")
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case compile_cmd.FullCommand():
			doCompile()

		default:
			return false
		}
		return true
	})
}
