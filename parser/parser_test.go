package parser

import (
	"fmt"
	"reflect"
	"testing"

	"bytes"
	"github.com/SHyx0rmZ/cgen/ast"
	goast "go/ast"
)

func TestParser_Parse(t *testing.T) {
	var tests []struct {
		Input string
		Value []ast.Node
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
