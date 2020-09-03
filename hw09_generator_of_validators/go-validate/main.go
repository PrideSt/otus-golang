package main

import (
	"flag"
	"fmt"
	ast "go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
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
	fOut, err := os.Create(outFile)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to open file %q for writing: %w", outFile, err))
	}
	defer func() {
		err := fOut.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("error when close file %q: %w", outFile, err))
		}
	}()

	if err := generator.Generate(fOut, pkgName, structs); err != nil {
		log.Fatal(err)
	}
}
