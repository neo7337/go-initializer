package router

import (
	"github.com/neo7337/go-initializer/handler"
	"oss.nandlabs.io/golly/rest"
)

func RouterHandler(server rest.Server) rest.Server {
	// Define your routing logic here

	server.Post("/v1/project", handler.CreateProject)

	return server
}
