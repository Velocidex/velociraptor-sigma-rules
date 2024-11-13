package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func indent(in string, indent int) string {
	indent_str := strings.Repeat(" ", indent)
	lines := strings.Split(in, "\n")
	result := []string{}
	for _, l := range lines {
		result = append(result, indent_str+l)
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

func ValuesInOrder[V interface{}](in map[string]V) (res []V) {
	keys := []string{}
	for k := range in {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		v, _ := in[k]
		res = append(res, v)
	}

	return res
}
