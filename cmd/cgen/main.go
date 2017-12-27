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
	}
}
