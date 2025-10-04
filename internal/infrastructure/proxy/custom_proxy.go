package proxy

import (
	"github.com/felisest/comproxy/internal/infrastructure/config"
	"go.uber.org/fx"
)

type CustomProxy struct {
	processor  func(...[]byte) error
	cfg        config.Configuration
	shutdowner fx.Shutdowner
}

func NewCustomProxy(
	processor func(...[]byte) error,
	cfg config.Configuration,
	shutdowner fx.Shutdowner,
) *CustomProxy {
	return &CustomProxy{
		cfg:       cfg,
		processor: processor,
	}
}
