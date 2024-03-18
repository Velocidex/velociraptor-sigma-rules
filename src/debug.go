package main

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func DebugPrint(msg string, args ...interface{}) {
	if *debug {
		fmt.Printf(msg, args...)
	}
}

func Debug(arg interface{}) {
	spew.Dump(arg)
}

func JsonDump(arg interface{}) string {
	serialized, _ := json.Marshal(arg)
	return string(serialized)
}
