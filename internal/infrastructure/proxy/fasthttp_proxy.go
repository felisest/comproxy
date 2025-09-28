package proxy

import (
	"context"

	"github.com/valyala/fasthttp"
)

type FastHttpProxy struct {
	CustomProxy
	ctx context.Context
}

func NewFastHttpProxy(
	ctx context.Context,
	proxy *CustomProxy,
) *FastHttpProxy {

	return &FastHttpProxy{
		ctx:         ctx,
		CustomProxy: *proxy,
	}
}

func (p *FastHttpProxy) Run() error {

	err := fasthttp.ListenAndServe(":"+p.cfg.Port, p.requestHandler)

	return err
}

func (p *FastHttpProxy) requestHandler(ctx *fasthttp.RequestCtx) {

	if string(ctx.Path()) == p.cfg.Path && ctx.IsPost() {

		body := ctx.Request.Body()

		request := make([]byte, len(body))
		copy(request, body)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		ctx.Request.CopyTo(req)
		req.SetRequestURI(p.cfg.RemoteHost + p.cfg.Path)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		client := &fasthttp.Client{}
		if err := client.Do(req, resp); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		resp.CopyTo(&ctx.Response)

		response := resp.Body()
		p.proc(request, response)

	} else {
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}
