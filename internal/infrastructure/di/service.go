package di

import (
	"context"

	"github.com/felisest/comproxy/internal/infrastructure/config"
	"github.com/felisest/comproxy/internal/infrastructure/logs"
	"github.com/felisest/comproxy/internal/infrastructure/proxy"
	"github.com/felisest/comproxy/internal/infrastructure/web"
	"github.com/felisest/comproxy/internal/operational/port"
	"github.com/felisest/comproxy/internal/operational/processor"

	"go.uber.org/fx"
)

func InitService() error {
	fxNew := fx.New(createApp())
	fxNew.Run()

	return fxNew.Err()
}

func createApp() fx.Option {
	return fx.Options(
		fx.Provide(
			context.Background,
			config.GetConfig,
		),
		fx.Provide(
			fx.Annotate(
				logs.NewZapLogger,
				fx.As(new(port.ILogger)),
			),
		),
		fx.Provide(
			proxy.NewCustomProxy,
			proxy.NewFastHttpProxy,
		),
		fx.Provide(
			processor.NewResponseComparer,
			func(p *processor.ResponseComparer) func([]byte, []byte) error {
				return p.GetProcedure()
			},
			fx.Annotate(
				web.NewHttpRequest,
				fx.As(new(port.IRequester)),
			),
		),
		fx.Invoke(
			func(p *proxy.FastHttpProxy) {
				p.Run()
			},
		),
	)
}
