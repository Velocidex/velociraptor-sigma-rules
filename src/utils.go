package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func indent(in string, indent int) string {
	lines := strings.Split(in, "\n")
	result := []string{}
	for _, l := range lines {
		result = append(result, strings.Repeat(" ", indent)+l)
	}
	return strings.Join(result, "\n")
}

func MustMarshal(in interface{}) []byte {
	data, _ := json.Marshal(in)
	return data
}

func MustMarshalIndent(in interface{}) []byte {
	data, _ := json.MarshalIndent(in, "", " ")
	return data
}

func Dump(in interface{}) {
	data, _ := json.MarshalIndent(in, "", " ")
	fmt.Printf("%v", string(data))
}
