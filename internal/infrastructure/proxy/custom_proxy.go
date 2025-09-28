package proxy

import (
	"github.com/felisest/comproxy/internal/infrastructure/config"
	"go.uber.org/fx"
)

type CustomProxy struct {
	proc func([]byte, []byte) error

	cfg *config.Config

	shutdowner fx.Shutdowner
}

func NewCustomProxy(
	proc func([]byte, []byte) error,
	cfg *config.Config,
	shutdowner fx.Shutdowner,
) *CustomProxy {
	return &CustomProxy{
		cfg:  cfg,
		proc: proc,
	}
}

func (p *CustomProxy) SetProcessor(proc func([]byte, []byte)) {
	//p.proc = proc
}
