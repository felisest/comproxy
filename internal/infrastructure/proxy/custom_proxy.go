package proxy

import (
	"github.com/felisest/comproxy/internal/infrastructure/config"
	"go.uber.org/fx"
)

type CustomProxy struct {
	Proc       func([]byte, []byte) error
	cfg        config.Configuration
	shutdowner fx.Shutdowner
}

func NewCustomProxy(
	proc func([]byte, []byte) error,
	cfg config.Configuration,
	shutdowner fx.Shutdowner,
) *CustomProxy {
	return &CustomProxy{
		cfg:  cfg,
		Proc: proc,
	}
}

func (p *CustomProxy) SetProcessor(proc func([]byte, []byte)) {
	//p.proc = proc
}
