package processor

import (
	"github.com/felisest/comproxy/internal/infrastructure/config"
	counter "github.com/felisest/comproxy/internal/infrastructure/counter"
	"github.com/felisest/comproxy/internal/operational/port"
)

type ResponseComparer struct {
	counter   counter.AtomicCounter
	requester port.IRequester
	logs      port.ILogger
	cfg       config.Configuration
	comparer  port.IComparer
}

func NewResponseComparer(
	cfg config.Configuration,
	requester port.IRequester,
	logs port.ILogger,
	comparer port.IComparer,
) *ResponseComparer {
	return &ResponseComparer{
		counter:   *counter.NewEventCounter(),
		requester: requester,
		cfg:       cfg,
		logs:      logs,
		comparer:  comparer,
	}
}

// [0] Request, [1] Response by remote
func (p *ResponseComparer) Process(data ...[]byte) error {
	p.counter.Inc()

	if p.counter.Value() > p.cfg.Proxy.Rate-1 {
		p.counter.Reset()

		//Request remote tested
		responseRemote, err := p.requester.Post(data[0])
		if err != nil {
			p.logs.Error("[PROCESS] Request error: ", err)
			return err
		}

		//Compare results
		ok, diff := p.comparer.Compare(data[1], responseRemote)

		if !ok {
			p.logs.Warn("Differ: %s", diff)
		}
	}

	return nil
}

func (p *ResponseComparer) GetProcedure() func(...[]byte) error {
	return p.Process
}
