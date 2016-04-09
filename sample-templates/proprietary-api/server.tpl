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
	    {{if $methodDef.HasQueryArguments}}
	        q := c.Query()
	    {{end}}
	    {{range $param := $methodDef.Parameters}}
	        {{if $param.InBody}}
                {{$param.GoName}} := api.{{$param.Schema.GoTypeName}}{}
                if !c.ParseRequest(&{{$param.GoName}}) {
                    return
            }{{end}}
            {{if $param.InPath}}
                {{$param.GoName}} := c.PathParam("{{$param.Name}}")
            {{end}}
            {{if $param.InQuery}}
                // type {{$param.Type}}
                {{if eq $param.Type "integer"}}
                    {{$param.GoName}} := q.IntValue("{{$param.Name}}", {{if $param.Default}}{{$param.Default}}{{else}}0{{end}})
                {{else if eq $param.Type "string"}}
                    {{$param.GoName}} := q.StringValue("{{$param.Name}}", "{{$param.Default}}")
                {{end}}


            {{end}}
            {{if $param.InHeader}}
                {{$param.GoName}} := c.Req.Header.Get("{{$param.Name}}")
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