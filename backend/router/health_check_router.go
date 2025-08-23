package router

import (
	"github.com/neo7337/go-initializer/response"
	"oss.nandlabs.io/golly/rest"
)

func HealthcheckRouterHandler(server rest.Server) rest.Server {
	// Health check endpoint.
	server.Get("/v1/healthz", func(ctx rest.ServerContext) {
		response.ResponseJSON(ctx, 200, map[string]string{"status": "ok"})
	})

	return server
}
