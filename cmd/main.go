package main

import (
	"github.com/bsonger/devflow-verify-service/pkg/router"
	"github.com/bsonger/devflow-verify-service/platform/shared/bootstrap"
)

func main() {
	err := bootstrap.Run(bootstrap.Options{
		Name: "verify-service",
		RouteOptions: router.Options{
			ServiceName:   "verify-service",
			EnableSwagger: true,
			Modules: []router.Module{
				router.ModuleVerify,
			},
		},
		PortEnv:        "VERIFY_SERVICE_PORT",
		DefaultPort:    8084,
		MetricsPortEnv: "VERIFY_SERVICE_METRICS_PORT",
		PprofPortEnv:   "VERIFY_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
