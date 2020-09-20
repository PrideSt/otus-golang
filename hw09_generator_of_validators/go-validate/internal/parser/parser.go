package parser

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

const (
	FieldKindType  = "type"
	FieldKindArray = "array"
	FieldKindClass = "class"
)

type StructDesc struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name         string
	DType        string
	InternalType string
	Kind         string
	Validators   []ValidatorDesc
}

type ValidatorDesc struct {
	FuncName string
	Args     []string
}

// ParseStructs analyze all type declarations, find structs among them and return map of structs with
// there fields description. Returns only types and fields which has validator tags.
func ParseStructs(f ast.Node) map[string]StructDesc {
	validators := make(map[string]StructDesc)

	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				if fields := getStructFields(structType); len(fields) > 0 {
					validators[typeSpec.Name.Name] = StructDesc{
						Name:   typeSpec.Name.Name,
						Fields: fields,
					}
				}
			}
		}

		return true
	})

	return validators
}

// getStructFields returns fields description of struct t, them which has some validator tags.
func getStructFields(t *ast.StructType) (result []Field) {
	for _, astField := range t.Fields.List {
		if len(astField.Names) > 0 {
			for _, fieldName := range astField.Names {
				fields := getType(fieldName.Name, astField.Type, getValidators(astField.Tag))
				result = append(result, fields...)
			}

			continue
		}
		// embedded types
		if indent, ok := astField.Type.(*ast.Ident); ok {
			if typeSpec, ok := indent.Obj.Decl.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					fields := getStructFields(structType)
					result = append(result, fields...)
				}
			}
		}
	}

	return result
}

// getType returns field type description: dynamic type- dType, internalType and kind.
func getType(name string, expr ast.Expr, v []ValidatorDesc) []Field {
	switch t := expr.(type) {
	case *ast.ArrayType:
		if len(v) > 0 {
			fields := getType(name, t.Elt, v)
			for i := range fields {
				if fields[i].Kind == FieldKindArray {
					log.Fatal(fmt.Errorf("nested arrays is not unsupported %+v", fields[i]))
				}
				fields[i].Kind = FieldKindArray
			}

			return fields
		}
	case *ast.Ident:
		if t.Obj != nil {
			fields := getRefType(name, t.Obj, v)
			for i := range fields {
				fields[i].DType = t.Name
			}
			return fields
		}

		// don't save scalar type without validators
		if len(v) > 0 {
			return []Field{{
				Name:         name,
				DType:        t.Name,
				InternalType: t.Name,
				Kind:         FieldKindType,
				Validators:   v,
			}}
		}

	case *ast.StructType:
		return getStructFields(t)
	default:
		log.Fatal(fmt.Errorf("unsuported ast.Expr type"))
	}
	return nil
}

// getRefType return fields from resource by reference.
func getRefType(name string, obj *ast.Object, v []ValidatorDesc) []Field {
	typeSpec, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return []Field{{
			Name:         name,
			DType:        obj.Name,
			InternalType: obj.Name,
		}}
	}

	var kind string
	switch typeSpec.Type.(type) {
	// reference to scalar type
	case *ast.Ident:
		kind = FieldKindType
	// reference to struct
	case *ast.StructType:
		kind = FieldKindClass
	default:
		log.Fatalf("struct subtype not supported for %q", name)
	}

	fields := getType(name, typeSpec.Type, v)
	for i := range fields {
		fields[i].Name = name
		fields[i].Kind = kind
	}

	return fields
}

func getValidators(tag *ast.BasicLit) []ValidatorDesc {
	if tag == nil {
		return nil
	}

	tags := parseTags(strings.Trim(tag.Value, " `'\""))

	return scanValidators(tags["validate"])
}

// parseTags lookup tags in given string s like `json:"foo,omitempty,string" xml:"foo"` and
// return slice of package -> tag, {"json" -> "foo,omitempty,string", "xml" -> "foo"} for given
// sample.
func parseTags(s string) map[string]string {
	result := make(map[string]string)
	for _, ss := range strings.Split(s, " ") {
		kvPair := strings.SplitN(ss, ":", 2)
		if len(kvPair) == 2 {
			if tags := strings.Trim(kvPair[1], " \"'"); len(tags) > 0 {
				result[kvPair[0]] = tags
			}
		}
	}
	return result
}

// scanValidators extracts validators from given string: for string `min:1|in:2,4,6` it returns
// slice of pair function name and them args like {{"min", []{1}}, {"in", []{2, 4, 6}}}.
func scanValidators(s string) (result []ValidatorDesc) {
	for _, ss := range strings.Split(s, "|") {
		if len(ss) == 0 {
			continue
		}

		kvPair := strings.SplitN(ss, ":", 2)
		if len(kvPair) < 2 {
			result = append(result, ValidatorDesc{
				FuncName: kvPair[0],
			})

			continue
		}

		result = append(result, ValidatorDesc{
			FuncName: kvPair[0],
			Args:     strings.Split(kvPair[1], ","),
		})
	}

	return result
}
