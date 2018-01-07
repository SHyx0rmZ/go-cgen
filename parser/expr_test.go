package parser

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"reflect"
	"testing"

	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/token"
)

func TestParser_ParseExpr(t *testing.T) {
	tests := []struct {
		Input string
		Value []ast.Node
	}{
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
