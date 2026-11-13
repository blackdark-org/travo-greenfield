package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/blackdark-org/travo-greenfield/internal/connectivity"
)

type connectivityOutput struct {
	Body connectivity.Status
}

func (s *server) registerConnectivity(a huma.API) {
	register(s, false, op("getDashboardConnectivity", http.MethodGet, "/api/v1/dashboard/connectivity", "Uplink and internet status", "dashboard"), a, s.getDashboardConnectivity)
}

func (s *server) getDashboardConnectivity(ctx context.Context, _ *struct{}) (*connectivityOutput, error) {
	status := connectivity.Status{Internet: connectivity.Unknown}
	if s.Connectivity != nil {
		status = s.Connectivity.Status(ctx)
	}
	return &connectivityOutput{Body: status}, nil
}
