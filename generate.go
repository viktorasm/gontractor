package gontractor

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"regexp"
)

type TagGeneratorFunc func (fieldName string, fieldDefinition SwaggerSchema, objectDefinition SwaggerSchema) string

type GeneratorSetup struct {
	tagGenerators [] TagGeneratorFunc
	buf bytes.Buffer
}

func (gs *GeneratorSetup) out(s string,params ...interface{}) {
	gs.buf.WriteString(fmt.Sprintf(s, params...))
}

func (s *GeneratorSetup) SetTagGenerators(tagGenerators ... TagGeneratorFunc) {
	s.tagGenerators = tagGenerators
}

func (s GeneratorSetup) generateTag(fieldName string, fieldDefinition SwaggerSchema, objectDefinition SwaggerSchema) string {
	result := make([]string,len(s.tagGenerators))
	for i,gen := range s.tagGenerators {
		result[i] = gen(fieldName,fieldDefinition,objectDefinition)
	}
	return "`"+strings.Join(result," ")+"`"
}

func JsonTags(fieldName string, fieldDefinition SwaggerSchema, objectDefinition SwaggerSchema) string {
	result := fieldName
	if !objectDefinition.IsRequired(fieldName) {
		result = result + ",omitempty"
	}
	return fmt.Sprintf("json:\"%s\"",result)
}

func (s GeneratorSetup) generateMethodName(httpMethod string, httpPath string, methodDef SwaggerPathOperation) string {
	if methodDef.OperationId!="" {
		return methodDef.OperationId
	}

	return httpMethod+"_"+strings.Trim(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(httpPath,"_"),"_")
}

func (gs *GeneratorSetup) writeComment(s string) {
	if s != "" {
		gs.out("// ")
		gs.out(s)
		gs.out("\n")
	}
}

func (gs *GeneratorSetup) writeSchemaDef(f SwaggerFile, s SwaggerSchema) {
	out := gs.out

	if s.Ref != "" {
		ref, err := f.findRefSchema(s.Ref)
		if err != nil {
			panic(err.Error())
		}
		out(ref.goTypeName)
		return
	}

	switch s.Type {
	case "array":
		out("[]")
		gs.writeSchemaDef(f,*s.Items)
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
			gs.writeComment(prop.Description)
			out("    ")
			out(strings.Title(name))
			out(" ")
			gs.writeSchemaDef(f, *prop)
			out(" ")
			out(gs.generateTag(name,*prop, s))

			out("\n")

		}
		out("}\n")

	default:
		panic("unknown schema type " + s.Type)
	}
}

func (s * GeneratorSetup) writeMethodParameter(f SwaggerFile, param SwaggerParameter) {
	out := s.out

	out(param.goName())
	out(" ")


	if param.Schema!=nil {
		out(param.Schema.goTypeName)
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
			panic("can't handle parameters of type "+param.Type)
		}
	}
}

func generateInterface(f SwaggerFile, generatorSetup GeneratorSetup) string {
	out := generatorSetup.out

	for _, definition := range f.Definitions {
		generatorSetup.writeComment(definition.Description)
		out("type %v ", definition.goTypeName)
		generatorSetup.writeSchemaDef(f,*definition)
		out("\n")
	}


	// interface
	out("type Service interface {\n")
	for httpPath,methodDefs := range f.Paths {
		for httpMethod, methodDef := range methodDefs {
			out("// %s %s\n",strings.ToUpper(httpMethod), httpPath)
			generatorSetup.writeComment(methodDef.Description)

			out(generatorSetup.generateMethodName(httpMethod,httpPath,*methodDef))
			out("(")
			for index, param := range methodDef.Parameters {
				if index!=0 {
					out(", ")
				}
				generatorSetup.writeMethodParameter(f, *param)
			}

			out(") (")

			for _,response := range methodDef.Responses {
				if response.Schema==nil {
					continue
				}
				out(response.Schema.goTypeName)
				out(", ")
			}

			out("error)")

			out("\n\n")
		}
	}
	out("}\n")


	//fmt.Println(generatorSetup.buf.String())
	formatted, err := format.Source(generatorSetup.buf.Bytes())
	if err != nil {
		panic("failed to format: " + err.Error())
	}



	return string(formatted)
}
