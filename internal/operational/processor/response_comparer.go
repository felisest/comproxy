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

// [0] Request, [1] Response by remote
func (p *ResponseComparer) Process(data ...[]byte) error {
	p.counter.Inc()

	if p.counter.Value() > p.cfg.Proxy.Rate-1 {
		p.counter.Reset()

		//Request remote tested
		responseRemote, _ := p.requester.Post(data[0])

		//Compare results
		ok, diff := CompareJson(data[1], responseRemote)

		if !ok {
			p.logs.Warn("Differ: %s", diff)
		}
	}

	return nil
}

func (p *ResponseComparer) GetProcedure() func(...[]byte) error {
	return p.Process
}
