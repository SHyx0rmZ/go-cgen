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

func TestParser_ParseDecl(t *testing.T) {
	tests := []struct {
		Input string
		Value []ast.Node
	}{
		{
			Input: "extern",
			Value: []ast.Node{
				&ast.ExternDecl{
					KeyPos: 0,
					Decl:   nil,
				},
			},
		},
		{
			Input: `extern "C" {`,
			Value: []ast.Node{
				&ast.ExternDecl{
					KeyPos: 0,
					Decl: &ast.CDecl{
						Value: &ast.BasicLit{
							ValuePos: 7,
							Kind:     token.STRING,
							Value:    `"C"`,
						},
						BodyPos: 11,
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
