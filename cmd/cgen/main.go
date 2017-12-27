package main

import (
	"bytes"
	"fmt"
	"github.com/SHyx0rmZ/cgen"
	"io"
	"os"
	"path/filepath"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b := new(bytes.Buffer)
	_, err = io.Copy(b, f)
	if err != nil {
		panic(err)
	}

	for i := range cgen.NewParser(filepath.Base(os.Args[1]), b.String()).Parse() {
		fmt.Printf("%#v\n", i)
		if td, ok := i.(cgen.TypeDecl); ok {
			if ed, ok := td.Expr.(cgen.EnumDecl); ok {
				fmt.Printf("typedef enum {\n")
				for _, s := range ed.Specs {
					switch es := s.(type) {
					case cgen.EnumValue:
						fmt.Printf("\t%s", es.Name.Name)
						if es.Value != nil {
							fmt.Printf(" = %s", es.Value.Value)
						}
						fmt.Printf(",\n")
					case cgen.EnumConstExpr:
						if be, ok := es.Expr.(cgen.BinaryExpr); ok {
							fmt.Printf("\t%s = %s | %s,\n", es.Name.Name, be.X.(cgen.Ident).Name, be.Y.(cgen.BasicLit).Value)
						}
					}
				}
				fmt.Printf("} %s;", td.Name.Name)
			}
		}
	}
}
