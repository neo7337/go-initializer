package main

import (
	"github.com/neo7337/go-initializer/router"
	"oss.nandlabs.io/golly/lifecycle"
	"oss.nandlabs.io/golly/rest"
)

var componentManager lifecycle.ComponentManager

func getRestServer(manager lifecycle.ComponentManager) lifecycle.Component {
	// Load application configuration
	conf := GetConfig()
	serverOptions := rest.DefaultSrvOptions()
	serverOptions.PathPrefix = "/api"
	serverOptions.ListenHost = conf.System.Host
	serverOptions.ListenPort = int16(conf.System.Port)
	serverOptions.ReadTimeout = int64(conf.System.ReadTimeout)   // seconds
	serverOptions.WriteTimeout = int64(conf.System.WriteTimeout) // seconds
	componentManager = manager

	server, err := rest.NewServer(serverOptions)
	if err != nil {
		panic(err)
	}
	server.OnChange(func(prevState, newState lifecycle.ComponentState) {
		if prevState == lifecycle.Unknown && newState == lifecycle.Starting {
			server = router.HealthcheckRouterHandler(server)
			server = router.RouterHandler(server)
		}
	})
	return server
}
