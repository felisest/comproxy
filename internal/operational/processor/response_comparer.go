package processor

import (
	"github.com/felisest/comproxy/internal/infrastructure/config"
	counter "github.com/felisest/comproxy/internal/operational/counter"
	"github.com/felisest/comproxy/internal/operational/port"
)

type ResponseComparer struct {
	counter   counter.AtomicCounter
	requester port.IRequester
	logs      port.ILogger
	cfg       config.Configuration
}

func NewResponseComparer(
	cfg config.Configuration,
	requester port.IRequester,
	logs port.ILogger,
) *ResponseComparer {
	return &ResponseComparer{
		counter:   *counter.NewEventCounter(),
		requester: requester,
		cfg:       cfg,
		logs:      logs,
	}
}

func (p *ResponseComparer) Process(request []byte, response []byte) error {
	p.counter.Inc()

	if p.counter.Value() > p.cfg.Proxy.Rate-1 {
		p.counter.Reset()

		//Request remote tested
		responseRemote, _ := p.requester.Post(request)

		//Compare results
		ok, diff := CompareJson(response, responseRemote)

		if !ok {
			p.logs.Warn("Differ: %s", diff)
		}

	}

	return nil
}

func (p *ResponseComparer) GetProcedure() func([]byte, []byte) error {
	return p.Process
}
