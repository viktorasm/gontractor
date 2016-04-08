package swagger

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamGoName(t *testing.T) {
	p := SwaggerParameter{
		Name: "voter-id",
	}

	require.Equal(t, "voterId", p.GoName())

}
