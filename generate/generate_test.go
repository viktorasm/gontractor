package generate

import (
	"fmt"
	"github.com/viktorasm/gontractor/swagger"
	"testing"
)

func TestSwaggerGen(t *testing.T) {
	result := swagger.Parse("../test-resources/swagger.yaml")
	//pretty.PrettyPrint(result)

	g := Generator{}
	g.SetTagGenerators(JsonTags)

	formattedInterface := g.GenerateApiInterface("api",*result)

	g = Generator{}
	g.SetTagGenerators(JsonTags)

	generatedServer := g.GenerateServerFromTemplate(*result, "../sample-templates/proprietary-api/server.tpl")
	fmt.Println("------------------- Interface ---")
	fmt.Println(formattedInterface)
	fmt.Println("------------------- Server ------")
	fmt.Println(generatedServer)
}
