package util

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetFiberCtx(fiber *fiber.App, host string, uri string, method string, body []byte) (*fiber.Ctx, func(), error) {
	var fReqCtx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	req.Header.SetHost(host)
	req.Header.SetContentType("application/json")
	_, err := req.BodyWriter().Write(body)
	if err != nil {
		return nil, nil, err
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		return nil, nil, err
	}
	fReqCtx.Init(&req, remoteAddr, nil)
	fctx := fiber.AcquireCtx(&fReqCtx)

	return fctx, func() {
		fiber.ReleaseCtx(fctx)
	}, nil
}
