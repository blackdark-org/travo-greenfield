// Package api registers every HTTP route; nothing is routed outside New.
package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type healthOutput struct {
	Body struct {
		Status  string `json:"status" enum:"ok" doc:"Always ok when the daemon answers"`
		Version string `json:"version" doc:"Running build version"`
	}
}

// New builds the router and its generated OpenAPI 3.1 spec.
func New(version string) (http.Handler, huma.API) {
	mux := http.NewServeMux()
	cfg := huma.DefaultConfig("travo-greenfield", version)
	cfg.OpenAPIPath = "/api/openapi"
	cfg.DocsPath = ""
	cfg.SchemasPath = ""
	cfg.CreateHooks = nil // no  links: that route does not exist
	a := humago.New(mux, cfg)

	huma.Register(a, huma.Operation{
		OperationID: "getHealth",
		Method:      http.MethodGet,
		Path:        "/api/health",
		Summary:     "Liveness check",
		Tags:        []string{"system"},
	}, func(context.Context, *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		out.Body.Version = version
		return out, nil
	})
	return mux, a
}
