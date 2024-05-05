package main

import "sort"

type RejectedRule struct {
	Path  string `json:"Path,omitempty"`
	Error string `json:"Error,omitempty"`
}

type Rejected struct {
	Rejects []RejectedRule `json:"Rejects"`
}

func (self RejectedRule) String() string {
	return self.Path + self.Error
}

func (self *CompilerContext) GetRejected() *Rejected {
	result := &Rejected{}

	for error_message, paths := range self.errored_rules {
		for _, p := range paths {
			result.Rejects = append(result.Rejects, RejectedRule{
				Error: error_message,
				Path:  p,
			})
		}
	}

	sort.Slice(result.Rejects, func(i, j int) bool {
		return result.Rejects[i].Path < result.Rejects[j].Path
	})

	return result
}
