package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFieldMappingLambda(t *testing.T) {
	cases := []struct {
		description string
		lambda      string
		// Expected offending event-root field name, or "" when the
		// mapping is well formed.
		want string
	}{
		// The bug this guard exists to catch (config/linux_ebpf_base.yaml).
		{"bare System (the LogonId bug)", "x=>System.UserID", "System"},
		{"bare EventData", "x=>EventData.argv", "EventData"},
		{"bare root inside a function argument", "x=>foo(System.X)", "System"},

		// Well-formed mappings taken from the real config - none flagged.
		{"System through the argument", "x=>x.System.UserID", ""},
		{"EventData through the argument", "x=>x.EventData.pathname", ""},
		{"function call, no event root", "x=>timestamp(epoch=now())", ""},
		{"event root nested in a function argument",
			"x=>process_tracker_get(id=x.System.ProcessID).Data.Username", ""},
		{"comparison expression", "x=>x.System.EventName = 'security_socket_connect'", ""},
		{"nested function with quoted argument",
			"x=>stat(filename=format(format='/proc/%v/cwd', args=x.System.ProcessID)).Data.Link", ""},
		{"root name appearing only inside a string literal", "x=>'System error'", ""},
	}

	for _, tc := range cases {
		got := CheckFieldMappingLambda(tc.lambda)
		if tc.want == "" {
			assert.Empty(t, got, "%v: %v", tc.description, tc.lambda)
		} else {
			assert.Equal(t, tc.want, got, "%v: %v", tc.description, tc.lambda)
		}
	}
}
