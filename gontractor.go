package main

import (
	"flag"
	"fmt"
	"github.com/viktorasm/gontractor/generate"
	"github.com/viktorasm/gontractor/swagger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Gontractor struct {
	spec           string
	serverTemplate string
	serverOutFile  string
	clientOutFile  string
	apiOutFile     string
}

func NewGontractor() *Gontractor {
	return &Gontractor{
		spec:           "swagger.yaml",
		serverTemplate: "server.tpl",
		serverOutFile:  "server_generated.go",
		clientOutFile:  "client/client.go",
		apiOutFile:     "api/api.go",
	}
}

func (g Gontractor) saveFile(fileName string, contents string) error {
	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, os.ModePerm)
	ioutil.WriteFile(fileName, []byte(contents), 0700)
	return nil
}

// locates workspace root in the fileName, and returns subdir from {workspace/src}
func (g Gontractor) getAbsolutePackagePath(fileName string) (string, error) {
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	// pick up everything after last "src/"
	result := filepath.Dir(abs[strings.LastIndex(abs, filepath.FromSlash("/src/"))+5:])

	// make sure we have correct slashes: /
	result = filepath.ToSlash(result)
	return result, nil
}

// guesses package name for given output Go file. handles relative urls
func (g Gontractor) getPackageName(fileName string) (string, error) {
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	return filepath.Base(filepath.Dir(abs)), nil
}

func (g Gontractor) Execute() error {
	spec := swagger.Parse(g.spec)
	generator := generate.Generator{}
	generator.SetTagGenerators(generate.JsonTags)

	apiContents, err := generator.GenerateApiInterface("api", *spec)
	if err != nil {
		return err
	}
	g.saveFile(g.apiOutFile, apiContents)

	templateData := generate.TemplateData{}
	templateData.Package.Api, err = g.getAbsolutePackagePath(g.apiOutFile)
	if err != nil {
		return err
	}

	templateData.Package.This, err = g.getPackageName(g.serverOutFile)
	if err != nil {
		return err
	}

	serverContents := generator.GenerateServerFromTemplate(*spec, g.serverTemplate, templateData)
	err = g.saveFile(g.serverOutFile, serverContents)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	g := NewGontractor()

	flag.StringVar(&g.spec, "spec", "swagger.yaml", "service specification flag")
	flag.StringVar(&g.serverTemplate, "server-template", "", "template to generate server")
	flag.Parse()

	err := g.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %s", err.Error())
		flag.Usage()
		os.Exit(1)
	}
}
