package cgen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/SHyx0rmZ/cgen/token"
)

func TestParser_ParsePreprocessorDirectives(t *testing.T) {
	tests := []struct {
		Input string
		Value []Node
	}{
		{
			"#define VALUE",
			[]Node{
				&MacroDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args:  nil,
					Value: nil,
				},
			},
		},
		{
			"#define VALUE()",
			[]Node{
				&MacroDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: &ArgList{
						Opening: 0,
						List:    nil,
						Closing: 0,
					},
					Value: nil,
				},
			},
		},
		{
			"#define VALUE 1",
			[]Node{
				&MacroDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: nil,
					Value: &BasicLit{
						ValuePos: 14,
						Kind:     token.INT,
						Value:    "1",
					},
				},
			},
		},
		{
			"#define VALUE(X) -1 / X",
			[]Node{
				&MacroDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: &ArgList{
						Opening: 0,
						List: []*Ident{
							{
								NamePos: 0,
								Name:    "X",
							},
						},
						Closing: 0,
					},
					Value: &BinaryExpr{
						X: &UnaryExpr{
							OpPos: 0,
							Op:    token.SUB,
							X: &BasicLit{
								ValuePos: 0,
								Kind:     token.INT,
								Value:    "1",
							},
						},
						OpPos: 0,
						Op:    token.QUO,
						Y: &Ident{
							NamePos: 0,
							Name:    "X",
						},
					},
				},
			},
		},
		{
			"#define VALUE (X) -1 / X",
			[]Node{
				&MacroDir{
					DirPos: 0,
					Name: &Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Value: &BinaryExpr{
						X: &NestedExpr{
							Opening: 0,
							Expr: &Ident{
								NamePos: 0,
								Name:    "X",
							},
							Closing: 0,
						},
						OpPos: 0,
						Op:    token.SUB,
						Y: &BinaryExpr{
							X: &BasicLit{
								ValuePos: 0,
								Kind:     token.INT,
								Value:    "1",
							},
							OpPos: 0,
							Op:    token.QUO,
							Y: &Ident{
								NamePos: 0,
								Name:    "X",
							},
						},
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
		{
			"#endif",
			[]Node{
				&EndIfDir{
					DirPos: 0,
				},
			},
		},
		// #endif
		// extern "C" {
		// #define (SD)
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
