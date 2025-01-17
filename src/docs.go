package main

import "time"

func (self *CompilerContext) GetDocs() (string, error) {
	params := &ArtifactContent{
		Time:       time.Now().UTC().Format(time.RFC3339),
		LogSources: BuildLogSource(self.config_obj),
	}

	templ, err := calculateTemplate(self.config_obj.DocTemplate, params)
	if err != nil {
		return "", err
	}

	return templ, nil
}
