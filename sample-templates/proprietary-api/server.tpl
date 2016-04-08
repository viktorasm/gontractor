package {{.Package.This}}

import (
	"github.com/gorilla/mux"
	"net/http"
	api "{{.Package.Api}}"
	servicecore_http "nevercrashed.com/servicecore/http"
)

func makeRoutes(s api.Service) http.Handler {
	m := mux.NewRouter()
	r := servicecore_http.NewRouteWrapper(m)

    {{range $httpPath, $httpMethods := .Spec.Paths}}
        {{range $httpMethod, $methodDef := $httpMethods}}
	r.{{ title $httpMethod }}("{{ $httpPath }}", func(c *servicecore_http.HTTPHandlerContext) {
	    {{range $param := $methodDef.Parameters}}
	        {{if eq $param.In "body"}}
                {{$param.GoName}} := api.{{$param.Schema.GoTypeName}}{}
                if !c.ParseRequest(&{{$param.GoName}}) {
                    return
            }{{end}}
            {{if eq $param.In "path"}}
                {{$param.GoName}} := c.PathParam("{{$param.Name}}")
            {{end}}
		{{end}}

		result,err := s.{{$methodDef.MethodCallSignature}}

		if err != nil {
			c.HandleError(err)
			return
		}
		c.Status({{$methodDef.SuccessHttpCode}})
		c.Result(result)
	})
    	{{end}}
	{{end}}

	return m
}