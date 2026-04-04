package main

import (
	"github.com/bsonger/devflow-service-common/bootstrap"
	"github.com/bsonger/devflow-service-common/observability"
	"github.com/bsonger/devflow-verify-service/pkg/config"
	"github.com/bsonger/devflow-verify-service/pkg/router"
)

func main() {
	err := bootstrap.Run(bootstrap.Options[config.Config, router.Options, string]{
		Name: "verify-service",
		RouteOptions: router.Options{
			ServiceName:   "verify-service",
			EnableSwagger: true,
			Modules: []router.Module{
				router.ModuleVerify,
			},
		},
		Load:        config.Load,
		InitRuntime: config.InitRuntime,
		NewRouter: func(opts router.Options) bootstrap.Runner {
			return router.NewRouterWithOptions(opts)
		},
		ResolveConfigPort: func(cfg *config.Config) int {
			if cfg != nil && cfg.Server != nil {
				return cfg.Server.Port
			}
			return 0
		},
		StartMetricsServer: observability.StartMetricsServer,
		StartPprofServer:   observability.StartPprofServer,
		PortEnv:            "VERIFY_SERVICE_PORT",
		DefaultPort:        8084,
		MetricsPortEnv:     "VERIFY_SERVICE_METRICS_PORT",
		PprofPortEnv:       "VERIFY_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
