package main

import (
	"bytes"
	"flag"
	"fmt"
	ast "go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path"

	"github.com/PrideSt/otus-golang/hw09_generator_of_validators/go-validate/internal/generator"
	"github.com/PrideSt/otus-golang/hw09_generator_of_validators/go-validate/internal/parser"
)

var (
	inFile string
)

func init() {
	flag.StringVar(&inFile, "in", "", "path to file with models to validators generate")
}

func main() {
	flag.Parse()

	if inFile == "" {
		log.Fatal(fmt.Errorf("invalid input path, please set falue for -in argument"))
	}

	fileContent, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal(fmt.Errorf("read file failed: %w", err))
	}

	fSet := token.NewFileSet()
	f, err := ast.ParseFile(fSet, "", fileContent, 0)
	if err != nil {
		log.Fatal(fmt.Errorf("parsing file failed: %w", err))
	}

	pkgName := f.Name.Name
	structs := parser.ParseStructs(f)

	outFile := fmt.Sprintf("%s_validation_generated.go", path.Join(path.Dir(inFile), pkgName))

	// we can write in file direct, but when Generate failed with error we already create generated file
	// and when log error with os.Exit defer with close never called. It's fix gocritic exitAfterDefer error.
	var buffer bytes.Buffer
	if err := generator.Generate(&buffer, pkgName, structs); err != nil {
		log.Fatal(err)
	}

	// I want FileMode 644 not 600 or less, disable gosec linter
	if err := ioutil.WriteFile(outFile, buffer.Bytes(), 0644); err != nil { //nolint:gosec
		log.Fatal(err)
	}
}
