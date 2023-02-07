package web

import (
	"net"
	"net/http"
)

type HandlerFunc func(ctx Context)

var _ Server = &HttpServer{}

type HttpServer struct {
	*router
}

func NewHTTPServer() *HttpServer {
	return &HttpServer{
		NewRouter(),
	}
}

type Server interface {
	http.Handler
	Start(path string) error
	//AddRoute(method string,path string,handleFunc HandlerFunc)
}

func (h *HttpServer) serve(ctx Context) {

}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//构建Context后执行业务
	ctx := Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HttpServer) Start(path string) error {
	l, err := net.Listen("tcp", path)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

//路由注册挪移到route中
//func (h *HttpServer) AddRoute(method string, path string, handleFunc HandlerFunc) {
//
//}
