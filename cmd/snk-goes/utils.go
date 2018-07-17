package main

import (
	"strings"

	// external
	"github.com/k0kubun/pp"
)

func prettyPrint(msg interface{}) {
	pp.Println(msg)
}

// borrowed from https://github.com/hashicorp/serf/blob/master/command/agent/flag_slice_value.go

// AppendSliceValue implements the flag.Value interface and allows multiple
// calls to the same variable to append a list.
type AppendSliceValue []string

func (s *AppendSliceValue) String() string {
	return strings.Join(*s, ",")
}

//
// Set will add another argument value to slice.
//
func (s *AppendSliceValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}
	if ok := sliceExists(*s, value); !ok {
		*s = append(*s, value)
	}
	return nil
}

func sliceExists(strSlice []string, input string) bool {
	for _, entry := range strSlice {
		if input == entry {
			return true
		}
	}
	return false
}

func dedupSliceValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}