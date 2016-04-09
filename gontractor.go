package main
import (
	"flag"
	"github.com/viktorasm/gontractor/swagger"
	"github.com/viktorasm/gontractor/generate"
)



type Gontractor struct {
	spec string
	serverTemplate string
	serverOutFile string
	clientOutFile string
	apiOutFile string
	outDir string
}

func NewGontractor() *Gontractor {
	return &Gontractor {
		spec:"swagger.yaml",
		serverOutFile: "server_generated.go",
		clientOutFile: "client/client.go",
		apiOutFile: "api/api.go",
		outDir: ".",
	}
}

func (g Gontractor) saveFile(fileName string,contents string) error {
	return nil
}

func (g Gontractor) Execute() error {
	spec := swagger.Parse(g.spec)
	generator := generate.Generator{}
	generator.SetTagGenerators(generate.JsonTags)

	apiContents := generator.GenerateApiInterface(*spec)
	g.saveFile(g.apiOutFile,apiContents)
	return nil
}


func main() {
	g := NewGontractor()

	flag.StringVar(&g.spec,"spec","swagger.yaml","service specification flag")
	flag.StringVar(&g.serverTemplate, "server-template","","template to generate server")



}
