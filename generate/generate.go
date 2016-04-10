package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/viktorasm/gontractor/swagger"
	"go/format"
	"regexp"
	"strings"
	"text/template"
	"errors"
)

type TagGeneratorFunc func(fieldName string, fieldDefinition swagger.SwaggerSchema, objectDefinition swagger.SwaggerSchema) string

type Generator struct {
	tagGenerators []TagGeneratorFunc
	buf           bytes.Buffer
}

// simplification of adding buffer contents
func (g *Generator) out(s string, params ...interface{}) {
	g.buf.WriteString(fmt.Sprintf(s, params...))
}

func (g *Generator) SetTagGenerators(tagGenerators ...TagGeneratorFunc) {
	g.tagGenerators = tagGenerators
}

func (g Generator) generateTag(fieldName string, fieldDefinition swagger.SwaggerSchema, objectDefinition swagger.SwaggerSchema) string {
	result := make([]string, len(g.tagGenerators))
	for i, gen := range g.tagGenerators {
		result[i] = gen(fieldName, fieldDefinition, objectDefinition)
	}
	return "`" + strings.Join(result, " ") + "`"
}

func JsonTags(fieldName string, fieldDefinition swagger.SwaggerSchema, objectDefinition swagger.SwaggerSchema) string {
	result := fieldName
	if !objectDefinition.IsRequired(fieldName) {
		result = result + ",omitempty"
	}
	return fmt.Sprintf("json:\"%s\"", result)
}

func (g Generator) generateMethodName(httpMethod string, httpPath string, methodDef swagger.SwaggerPathOperation) string {
	if methodDef.OperationId != "" {
		return methodDef.OperationId
	}

	return httpMethod + "_" + strings.Trim(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(httpPath, "_"), "_")
}

func (g *Generator) writeComment(s string) {
	if s != "" {
		g.out("// ")
		g.out(s)
		g.out("\n")
	}
}

func (g *Generator) writeSchemaDef(f swagger.SwaggerSpec, s swagger.SwaggerSchema) error {
	out := g.out

	if s.Ref != "" {
		ref, err := f.FindRefSchema(s.Ref)
		if err != nil {
			return err
		}
		out(ref.GoTypeName)
		return nil
	}

	switch s.Type {
	case "array":
		out("[]")
		err := g.writeSchemaDef(f, *s.Items)
		if err!=nil{
			return err
		}
	case "boolean":
		out("bool")
	case "integer":
		out("int")
	case "number":
		out("float32")
	case "string":
		out("string")

	case "object":
		out("struct {\n")
		for name, prop := range *s.Properties {
			g.writeComment(prop.Description)
			out("    ")
			out(strings.Title(name))
			out(" ")
			err := g.writeSchemaDef(f, *prop)
			if err != nil {
				return err
			}
			out(" ")
			out(g.generateTag(name, *prop, s))

			out("\n")

		}
		out("}\n")

	default:
		return errors.New("unknown schema type " + s.Type)
	}

	return nil
}

func (g *Generator) writeMethodParameter(f swagger.SwaggerSpec, param swagger.SwaggerParameter) {
	out := g.out

	out(param.GoName())
	out(" ")

	if param.Schema != nil {
		out(param.Schema.GoTypeName)
	} else {
		switch param.Type {
		case "boolean":
			out("bool")
		case "integer":
			out("int")
		case "number":
			out("float32")
		case "string":
			out("string")
		case "object":

		default:
			panic("can't handle parameters of type " + param.Type)
		}
	}
}

// generates common API: request/response types (based on #/definitions/*), and
// a service interface (a function for for each HTTP method in each path)
func (g Generator) GenerateApiInterface(packageName string, f swagger.SwaggerSpec) (string,error) {
	out := g.out

	out("package %s\n", packageName)

	for _, definition := range f.Definitions {
		g.writeComment(definition.Description)
		out("type %v ", definition.GoTypeName)
		err := g.writeSchemaDef(f, *definition)
		if err!=nil{
			return "",err
		}
		out("\n")
	}

	// interface
	out("type Service interface {\n")
	for httpPath, methodDefs := range f.Paths {
		for httpMethod, methodDef := range methodDefs {
			out("// %s %s\n", strings.ToUpper(httpMethod), httpPath)
			g.writeComment(methodDef.Description)

			methodDef.GoInfo.InterfaceMethodName = strings.Title(g.generateMethodName(httpMethod, httpPath, *methodDef))
			out(methodDef.GoInfo.InterfaceMethodName)
			out("(")
			for index, param := range methodDef.Parameters {
				if index != 0 {
					out(", ")
				}
				g.writeMethodParameter(f, *param)
			}

			out(") (")

			for _, response := range methodDef.Responses {
				if response.Schema == nil {
					continue
				}
				out(response.Schema.GoTypeName)
				out(", ")
			}

			out("error)")

			out("\n\n")
		}
	}
	out("}\n")

	//fmt.Println(generatorSetup.buf.String())
	return g.formattedOutput(),nil
}

func (g Generator) formattedOutput() string {

	formatted, err := format.Source(g.buf.Bytes())
	if err != nil {
		println("failed to format:")
		return g.buf.String()
	}

	return string(formatted)
}

type TemplateData struct {
	Package struct {
		This string
		Api  string
	}
	Spec *swagger.SwaggerSpec
}

func (g Generator) GenerateServerFromTemplate(f swagger.SwaggerSpec, templateFileName string, templateData TemplateData) string {
	templateData.Spec = &f

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
	}

	t := template.Must(template.New("server.tpl").Funcs(funcMap).ParseFiles(templateFileName))
	w := bufio.NewWriter(&g.buf)
	err := t.Execute(w, &templateData)
	if err != nil {
		panic(err.Error())
	}
	w.Flush()
	return g.formattedOutput()
}
