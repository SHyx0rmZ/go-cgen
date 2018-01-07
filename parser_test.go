package cgen

import (
	"fmt"
	"reflect"
	"testing"

	"bytes"
	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/token"
	goast "go/ast"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		Input string
		Value []ast.Node
	}{
		{
			"#define VALUE",
			[]ast.Node{
				&ast.MacroDir{
					DirPos: 0,
					Name: &ast.Ident{
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
			[]ast.Node{
				&ast.MacroDir{
					DirPos: 0,
					Name: &ast.Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: &ast.ArgList{
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
			[]ast.Node{
				&ast.MacroDir{
					DirPos: 0,
					Name: &ast.Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: nil,
					Value: &ast.BasicLit{
						ValuePos: 14,
						Kind:     token.INT,
						Value:    "1",
					},
				},
			},
		},
		{
			"#define VALUE(X) -1 / X",
			[]ast.Node{
				&ast.MacroDir{
					DirPos: 0,
					Name: &ast.Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Args: &ast.ArgList{
						Opening: 13,
						List: []*ast.Ident{
							{
								NamePos: 14,
								Name:    "X",
							},
						},
						Closing: 15,
					},
					Value: &ast.BinaryExpr{
						X: &ast.UnaryExpr{
							OpPos: 17,
							Op:    token.SUB,
							X: &ast.BasicLit{
								ValuePos: 18,
								Kind:     token.INT,
								Value:    "1",
							},
						},
						OpPos: 20,
						Op:    token.QUO,
						Y: &ast.Ident{
							NamePos: 22,
							Name:    "X",
						},
					},
				},
			},
		},
		{
			"#define VALUE (X) -1 / X",
			[]ast.Node{
				&ast.MacroDir{
					DirPos: 0,
					Name: &ast.Ident{
						NamePos: 8,
						Name:    "VALUE",
					},
					Value: &ast.BinaryExpr{
						X: &ast.ParenExpr{
							Opening: 14,
							Expr: &ast.Ident{
								NamePos: 15,
								Name:    "X",
							},
							Closing: 16,
						},
						OpPos: 18,
						Op:    token.SUB,
						Y: &ast.BinaryExpr{
							X: &ast.BasicLit{
								ValuePos: 19,
								Kind:     token.INT,
								Value:    "1",
							},
							OpPos: 21,
							Op:    token.QUO,
							Y: &ast.Ident{
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
			[]ast.Node{
				&ast.IncludeDir{
					DirPos:  0,
					PathPos: 9,
					Path:    `"stddef.h"`,
				},
			},
		},
		{
			`#include <stddef.h>`,
			[]ast.Node{
				&ast.IncludeDir{
					DirPos:  0,
					PathPos: 9,
					Path:    `<stddef.h>`,
				},
			},
		},
		{
			"#endif",
			[]ast.Node{
				&ast.EndIfDir{
					DirPos: 0,
				},
			},
		},
		{
			"#else",
			[]ast.Node{
				&ast.ElseDir{
					DirPos: 0,
				},
			},
		},
		{
			"1",
			[]ast.Node{
				&ast.BasicLit{
					ValuePos: 0,
					Kind:     token.INT,
					Value:    "1",
				},
			},
		},
		{
			"-1",
			[]ast.Node{
				&ast.UnaryExpr{
					OpPos: 0,
					Op:    token.SUB,
					X: &ast.BasicLit{
						ValuePos: 1,
						Kind:     token.INT,
						Value:    "1",
					},
				},
			},
		},
		{
			"1 / 1",
			[]ast.Node{
				&ast.BinaryExpr{
					X: &ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    "1",
					},
					OpPos: 2,
					Op:    token.QUO,
					Y: &ast.BasicLit{
						ValuePos: 4,
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
				goast.Fprint(bufGot, nil, actual, goast.NotNilFilter)
				bufWant := new(bytes.Buffer)
				goast.Fprint(bufWant, nil, test.Value, goast.NotNilFilter)
				t.Errorf("%s:\ngot:\n%swant:\n%s", parser.name, bufGot.String(), bufWant.String())
			}
		})
	}
}
