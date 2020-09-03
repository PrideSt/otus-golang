package generator

import (
	"io"
	"text/template"

	"github.com/PrideSt/otus-golang/hw09_generator_of_validators/go-validate/internal/parser"
)

// Generate render validator function for types described in s and write Validate() method to w.
func Generate(w io.Writer, pkg string, s map[string]parser.StructDesc) error {
	mainTpl, err := template.New("validatorTpl").Funcs(funcs).Parse(fileTemplate)
	if err != nil {
		return err
	}

	for tName, tCont := range templates {
		if _, err := mainTpl.New(tName).Parse(tCont); err != nil {
			return err
		}
	}

	err = mainTpl.ExecuteTemplate(w, "validatorTpl", struct {
		Package string
		Imports []string
		Structs map[string]parser.StructDesc
	}{
		Package: pkg,
		Imports: []string{
			"fmt",
			"regexp",
			"",
			"github.com/PrideSt/otus-golang/hw09_generator_of_validators/validator",
		},
		Structs: s,
	})
	if err != nil {
		return err
	}

	return nil
}
