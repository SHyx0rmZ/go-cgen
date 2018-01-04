package cgen

import (
	"fmt"
	"reflect"
	"testing"

	"bytes"
	"github.com/SHyx0rmZ/cgen/token"
	"go/ast"
)

func TestParser_Parse(t *testing.T) {
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
						Opening: 13,
						List:    nil,
						Closing: 14,
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
						Opening: 13,
						List: []*Ident{
							{
								NamePos: 14,
								Name:    "X",
							},
						},
						Closing: 15,
					},
					Value: &BinaryExpr{
						X: &UnaryExpr{
							OpPos: 17,
							Op:    token.SUB,
							X: &BasicLit{
								ValuePos: 18,
								Kind:     token.INT,
								Value:    "1",
							},
						},
						OpPos: 20,
						Op:    token.QUO,
						Y: &Ident{
							NamePos: 22,
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
						X: &ParenExpr{
							Opening: 14,
							Expr: &Ident{
								NamePos: 15,
								Name:    "X",
							},
							Closing: 16,
						},
						OpPos: 18,
						Op:    token.SUB,
						Y: &BinaryExpr{
							X: &BasicLit{
								ValuePos: 19,
								Kind:     token.INT,
								Value:    "1",
							},
							OpPos: 21,
							Op:    token.QUO,
							Y: &Ident{
								NamePos: 23,
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
		{
			"1",
			[]Node{
				&BasicLit{
					ValuePos: 0,
					Kind:     token.INT,
					Value:    "1",
				},
			},
		},
		{
			"-1",
			[]Node{
				&UnaryExpr{
					OpPos: 0,
					Op:    token.SUB,
					X: &BasicLit{
						ValuePos: 1,
						Kind:     token.INT,
						Value:    "1",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.Input), func(t *testing.T) {
			parser := NewParser(t.Name(), test.Input)
			actual := parser.Nodes()

			if !reflect.DeepEqual(actual, test.Value) {
				bufGot := new(bytes.Buffer)
				ast.Fprint(bufGot, nil, actual, ast.NotNilFilter)
				bufWant := new(bytes.Buffer)
				ast.Fprint(bufWant, nil, test.Value, ast.NotNilFilter)
				t.Errorf("%s:\ngot:\n%swant:\n%s", parser.name, bufGot.String(), bufWant.String())
			}
		})
	}
}
