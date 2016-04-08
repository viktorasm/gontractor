package gontractor

import (
	"testing"
	//	"github.com/gobs/pretty"
	"fmt"
	"github.com/stretchr/testify/require"
)

func TestSwaggerGen(t *testing.T) {
	result := swaggerGen("swagger.yaml")
	//pretty.PrettyPrint(result)

	opts := GeneratorSetup{}
	opts.SetTagGenerators(JsonTags)

	formattedInterface := generateInterface(*result, opts)
	fmt.Println("-------------------")
	fmt.Println(formattedInterface)
}

func TestParamGoName(t *testing.T) {
	p := SwaggerParameter{
		Name: "voter-id",
	}

	require.Equal(t, "voterId", p.goName())

}
