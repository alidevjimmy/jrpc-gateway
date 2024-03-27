package internal

import (
	"bytes"
	"fmt"
	"text/template"
	"unicode"
)

var tmplFuncs = map[string]interface{}{
	"rpcMethod": rpcMethod,
}

// FileTmpl is .jgw template
var FileTmpl = template.Must(template.New("").Funcs(tmplFuncs).Parse(`
// Code generated by protoc-gen-jrpc-gateway. DO NOT EDIT.
// source: {{ .Name }}

{{$package := .Package}}
/*

Package {{ $package }} is a reverse proxy.

It translates gRPC into JSON-RPC 2.0
*/
package {{ $package }}

import (
	"context"
	"encoding/json"

	"github.com/jrpc-gateway/jsonrpc"
	"google.golang.org/protobuf/protobuf/encoding/protojson"
)

{{range $service := .Service}}

{{$serviceName := $service.GetName | printf "%sJsonRpcService"}}
{{$clientName := $service.GetName | printf "%sClient"}}

type {{$serviceName}} struct {
	client	{{$clientName}}
}

func {{$serviceName | printf "New%s"}} (client {{$clientName}}) {{$serviceName}} {
	return {{$serviceName}} {
		client: client,
	}
}

func (s *{{$serviceName}}) Methods() map[string]jsonrpc.Method {
	return map[string]func(ctx context.Context, params json.RawMessage) (any, error) {
		{{range $method := $service.GetMethod}}
			"{{rpcMethod $package $service.GetName $method.GetName}}": func(ctx context.Context, data json.RawMessage) (any, error) {
				req := new({{$method.GetName | printf "%sRequest"}})
				err := protojson.Unmarshal(data, req)
				if err != nil {
					return nil, err
				}
				return s.client.{{$method.GetName}}(ctx, req)
			},
		{{end}}
	}
}

{{end}}
`))

func rpcMethod(pkg, service, method string) string {
	return fmt.Sprintf("%s.%s.%s", camelToSnake(pkg), camelToSnake(service), camelToSnake(method))
}

func camelToSnake(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}