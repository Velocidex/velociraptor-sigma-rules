package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Velocidex/yaml/v2"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	profile_cmd        = app.Command("gen_profiles", "Generate profile JSON")
	profile_cmd_output = profile_cmd.Flag("output", "File to write the profile").Required().String()
	profile_cmd_args   = profile_cmd.Arg("configs", "Config files to read").Required().Strings()
)

func doProfile() error {
	profiles := make(map[string]*Config)

	for _, config_file := range *profile_cmd_args {
		fd, err := os.Open(config_file)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}

		conf := NewConfig()
		err = yaml.Unmarshal(data, conf)
		if err != nil {
			return err
		}
		profiles[conf.Name] = conf
	}

	serialized, err := json.Marshal(profiles)
	if err != nil {
		return err
	}

	outfd, err := os.OpenFile(*profile_cmd_output,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outfd.Close()

	_, err = outfd.Write(serialized)
	return err
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case profile_cmd.FullCommand():
			err := doProfile()
			kingpin.FatalIfError(err, "Compiling profiles")

		default:
			return false
		}
		return true
	})
}
