package data

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func wireApp(logger log.Logger, cfg config.Config, hs *http.Server, gs *grpc.Server, data *Data, sr *registry.ServiceRegistry) (*kratos.App, func(), error) {
	cleanup := func() {
		data.Close()
	}

	return kratos.New(
			kratos.Name("feishu-office-suite"),
			kratos.Logger(logger),
			kratos.Server(
				hs,
				gs,
			),
			kratos.Registrar(sr),
		), cleanup, nil
}