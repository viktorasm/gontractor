package swagger

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestParamGoName(t *testing.T) {
	p := SwaggerParameter{
		Name: "voter-id",
	}

	require.Equal(t, "voterId", p.GoName())

}
