package processor

import (
	"fmt"

	"github.com/felisest/comproxy/internal/infrastructure/config"
	counter "github.com/felisest/comproxy/internal/operational/counter"
	"github.com/felisest/comproxy/internal/operational/port"
)

type ResponseComparer struct {
	counter   counter.AtomicCounter
	requester port.IRequester

	cfg config.Config
}

func NewResponseComparer(
	cfg *config.Config,
	requester port.IRequester,
) *ResponseComparer {
	return &ResponseComparer{
		counter:   *counter.NewEventCounter(),
		requester: requester,
		cfg:       *cfg,
	}
}

func (p *ResponseComparer) Process(request []byte, response []byte) error {

	p.counter.Inc()

	if p.counter.Value() > p.cfg.Rate-1 {
		p.counter.Reset()

		//Request remote tested
		responseRemote, _ := p.requester.Post(request)

		//Compare results
		ok, diff := CompareJson(response, responseRemote)

		if !ok {
			fmt.Println(diff)
		}

	}

	return nil
}

func (p *ResponseComparer) GetProcedure() func([]byte, []byte) error {
	return p.Process
}
