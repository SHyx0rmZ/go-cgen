package cgen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/SHyx0rmZ/cgen/token"
)

func TestParser_ParseDefine(t *testing.T) {
	tests := []struct {
		Input string
		Value []Node
	}{
		{
			"#define VALUE",
			[]Node{
				&BadDir{
					From: 0,
					To:   13,
				},
			},
		},
		{
			"#define VALUE 1",
			[]Node{
				&DefineDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Value: &BasicLit{
						ValuePos: 14,
						Kind:     token.INT,
						Value:    "1",
					},
				},
			},
		},
		{
			`#include "stddef.h"`,
			[]Node{
				&IncludeDir{
					DirPos:  0,
					PathPos: 9,
					Path:    `"stddef.h"`,
				},
			},
		},
		{
			`#include <stddef.h>`,
			[]Node{
				&IncludeDir{
					DirPos:  0,
					PathPos: 9,
					Path:    `<stddef.h>`,
				},
			},
		},
	}

	for i, test := range tests {
		parser := NewParser(fmt.Sprintf("test #%d", i), test.Input)
		actual := parser.Nodes()

		if !reflect.DeepEqual(actual, test.Value) {
			bs, _ := json.MarshalIndent(actual, "", "  ")
			t.Errorf("%s: got %s", parser.name, string(bs))
		}
	}
}
