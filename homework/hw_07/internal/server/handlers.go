package server

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func (s *Server) stdHttpHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
func (s *Server) fastHttpHealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetBody([]byte("OK"))
}
