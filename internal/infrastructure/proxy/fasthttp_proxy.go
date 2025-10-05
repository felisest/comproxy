package proxy

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/felisest/comproxy/internal/operational/port"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

const shutdownTimeout = 3

type FastHttpProxy struct {
	CustomProxy
	ctx        context.Context
	logs       port.ILogger
	shutdowner fx.Shutdowner
}

func NewFastHttpProxy(
	ctx context.Context,
	proxy *CustomProxy,
	logs port.ILogger,
	shutdowner fx.Shutdowner,
) *FastHttpProxy {
	return &FastHttpProxy{
		ctx:         ctx,
		CustomProxy: *proxy,
		logs:        logs,
		shutdowner:  shutdowner,
	}
}

func (p *FastHttpProxy) Run() {
	router := provideRoutes(p.requestHandler)

	server := fasthttp.Server{
		Handler:            router.Handler,
		ReadTimeout:        p.cfg.Proxy.Server.GetTimeout(),
		WriteTimeout:       p.cfg.Proxy.Server.GetTimeout(),
		DisableKeepalive:   true,
		TCPKeepalive:       false,
		MaxRequestsPerConn: 1,
	}

	errChan := make(chan error)
	go func() {
		p.logs.Info("[FASTHTTP] Proxy starting on: ", p.cfg.Proxy.Server.Port)
		errChan <- server.ListenAndServe(":" + p.cfg.Proxy.Server.Port)
	}()

	go func() {
		defer func() {
			_ = p.shutdowner.Shutdown()
		}()

		p.logs.Info("[FASTHTTP] Proxy is running...")

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-sigCh:
			p.logs.Info("[FASTHTTP] Shutting down. Received signal: ", sig)

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
			defer cancel()

			if err := server.ShutdownWithContext(ctx); err != nil {
				p.logs.Error("[FASTHTTP] Shutdown error: ", err)
			}

		case err := <-errChan:
			p.logs.Error("[FASTHTTP] ListenAndServe error: ", err)

		}
	}()
}

func (p *FastHttpProxy) requestHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) == p.cfg.Proxy.Path && ctx.IsPost() {
		requestBody := ctx.Request.Body()
		request := make([]byte, len(requestBody))
		copy(request, requestBody)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		ctx.Request.CopyTo(req)
		req.SetRequestURI(p.cfg.Proxy.RemoteHost + p.cfg.Proxy.Path)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		client := &fasthttp.Client{}
		if err := client.Do(req, resp); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		resp.CopyTo(&ctx.Response)

		responseBody := resp.Body()
		response := make([]byte, len(responseBody))
		copy(response, responseBody)
		p.processor(request, response)

	} else {
		ctx.Error("[FASTHTTP] Not found", fasthttp.StatusNotFound)
	}
}

func provideRoutes(h fasthttp.RequestHandler) *router.Router {
	r := router.New()

	r.PanicHandler = func(ctx *fasthttp.RequestCtx, i any) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBodyString("Internal server error")
	}

	// Liveness, Readiness probes

	r.ANY("/{path:*}", h)

	return r
}
