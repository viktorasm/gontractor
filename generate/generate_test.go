package generate
import (
	"testing"
	"github.com/viktorasm/gontractor/swagger"
	"fmt"
)


func TestSwaggerGen(t *testing.T) {
	result := swagger.Parse("../test-resources/swagger.yaml")
	//pretty.PrettyPrint(result)

	opts := GeneratorSetup{}
	opts.SetTagGenerators(JsonTags)

	formattedInterface := generateInterface(*result, opts)
	fmt.Println("-------------------")
	fmt.Println(formattedInterface)
}
