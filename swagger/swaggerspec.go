package swagger

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type SwaggerTypedObject struct {
	Schema *SwaggerSchema `yaml:"schema"`
}

type SwaggerParameter struct {
	SwaggerTypedObject `yaml:",inline"`
	Ref                string `yaml:"$ref"`
	Name               string `yaml:"name"`
	In                 string `yaml:"in"`
	Description        string `yaml:"description"`
	Required           bool   `yaml:"required"`
	Type               string `yaml:"type"`
}

func (p SwaggerParameter) GoName() string {
	re := regexp.MustCompile("[^a-zA-Z]")
	result := re.ReplaceAllString(p.Name, " ")
	result = strings.Title(result)
	result = re.ReplaceAllString(result, "")
	result = p.Name[0:1] + result[1:]
	return result
}

type SwaggerSchema struct {
	GoTypeName  string
	Ref         string                     `yaml:"$ref"`
	Type        string                     `yaml:"type"`
	ReadOnly    bool                       `yaml:"readOnly"`
	Properties  *map[string]*SwaggerSchema `yaml:"properties"`
	Description string                     `yaml:"description"`
	Items       *SwaggerSchema             `yaml:"items"`
	Required    []string                   `yaml:"required"`
}

func (s SwaggerSchema) IsRequired(field string) bool {
	if s.Properties == nil {
		return false
	}
	for f := range s.Required {
		if f == field {
			return true
		}

	}
	return false
}

type SwaggerPathOperation struct {
	GoInfo struct {
		InterfaceMethodName string
	} `yaml:"-"`

	OperationId string                                   `yaml:"operationId"`
	Description string                                   `yaml:"description"`
	Parameters  []*SwaggerParameter                      `yaml:"parameters"`
	Responses   map[string]*SwaggerPathOperationResponse `yaml:"responses"`
}

func (op SwaggerPathOperation) MethodCallSignature() string {
	result := op.GoInfo.InterfaceMethodName + "("
	for index, param := range op.Parameters {
		if index > 0 {
			result += ", "
		}
		result += param.GoName()
	}
	result += ")"
	return result
}

func (op SwaggerPathOperation) SuccessHttpCode() string {
	for key := range op.Responses {
		return key
	}
	return "http.StatusOK"
}

type SwaggerPathOperationResponse struct {
	SwaggerTypedObject `yaml:",inline"`
	Description        string `yaml:"description"`
}

type SwaggerFile struct {
	BasePath    string                                      `yaml:"basePath"`
	Paths       map[string]map[string]*SwaggerPathOperation `yaml:"paths"`
	Parameters  map[string]*SwaggerParameter                `yaml:"parameters"`
	Definitions map[string]*SwaggerSchema                   `yaml:"definitions"`
}

func loadFile(inputFile string) (*SwaggerFile, error) {
	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	result := &SwaggerFile{}
	err = yaml.Unmarshal(file, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f SwaggerFile) replaceReferences() error {

	updateParamSchema := func(p *SwaggerTypedObject) {
		if p.Schema != nil && p.Schema.Ref != "" {
			var err error
			p.Schema, err = f.FindRefSchema(p.Schema.Ref)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	var replaceSchemaReferences func(s *SwaggerSchema) error
	replaceSchemaReferences = func(s *SwaggerSchema) error {
		if s.Properties != nil {
			for _, prop := range *s.Properties {
				err := replaceSchemaReferences(prop)
				if err != nil {
					return err
				}
			}
		}
		if s.Items != nil && s.Items.Ref != "" {
			schema, err := f.FindRefSchema(s.Items.Ref)
			if err != nil {
				return err
			}
			s.Items = schema
		}
		return nil
	}

	for _, pathOperations := range f.Paths {
		for _, operationDescr := range pathOperations {
			for i := 0; i < len(operationDescr.Parameters); i++ {
				p := operationDescr.Parameters[i]
				ref := p.Ref
				if ref != "" {
					p, err := f.findRefParam(ref)
					if err != nil {
						return err
					}
					operationDescr.Parameters[i] = p
				}
				updateParamSchema(&operationDescr.Parameters[i].SwaggerTypedObject)
			}
			for _, resp := range operationDescr.Responses {
				updateParamSchema(&resp.SwaggerTypedObject)
			}
		}
	}
	return nil
}

func (f SwaggerFile) generateGoTypeNames() {
	for key, value := range f.Definitions {
		value.GoTypeName = strings.Title(key)
	}
}

func (f SwaggerFile) findRefParam(ref string) (*SwaggerParameter, error) {
	ref = strings.TrimPrefix(ref, "#/parameters/")
	result, ok := f.Parameters[ref]
	if !ok {
		return nil, fmt.Errorf("param `%v` is not defined in spec", ref)
	}
	return result, nil
}

func (f SwaggerFile) FindRefSchema(ref string) (*SwaggerSchema, error) {
	refName := strings.TrimPrefix(ref, "#/definitions/")
	result, ok := f.Definitions[refName]
	if !ok {
		return nil, fmt.Errorf("definition `%v` is not defined in spec", ref)
	}
	return result, nil
}

func Parse(inputfile string) *SwaggerFile {
	result, err := loadFile(inputfile)
	if err != nil {
		panic(err.Error())
	}
	err = result.replaceReferences()
	if err != nil {
		panic(err.Error())
	}

	result.generateGoTypeNames()

	return result
}