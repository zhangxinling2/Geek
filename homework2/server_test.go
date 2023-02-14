package web

//import (
//	"gitee.com/geektime-geekbang/geektime-go/web/homework2/middleware/accesslog"
//	"gitee.com/geektime-geekbang/geektime-go/web/homework2/middleware/errhdl"
//	"gitee.com/geektime-geekbang/geektime-go/web/homework2/middleware/opentelemetry"
//	"gitee.com/geektime-geekbang/geektime-go/web/homework2/middleware/prometheus"
//	"go.opentelemetry.io/otel"
//	"net/http"
//	"testing"
//)
//const instrument=""
//func TestHTTPServer_UseV1(t *testing.T) {
//	accMiddleware:=accesslog.NewBuilder().Build()
//	errhdlMiddleware:=errhdl.NewMiddlewareBuilder().Build()
//	tracer:=otel.GetTracerProvider().Tracer(instrument)
//	traceMiddleware:=&opentelemetry.MiddlewareBuilder{
//		Tracer:tracer,
//	}
//	prometheusMiddleware:=&prometheus.MiddlewareBuilder{
//		Subsystem: "web",
//		Name: "http_response",
//		ConstLabels: map[string]string{},
//		Help: "help",
//	}
//	server:=NewHTTPServer()
//	testRoutes := []struct {
//		method string
//		path   string
//		mdls []Middleware
//	}{
//		{
//			method: http.MethodGet,
//			path:   "/",
//			mdls: []Middleware{
//				accMiddleware,
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user",
//			mdls: []Middleware{
//				errhdlMiddleware,
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user/home",
//			mdls: []Middleware{
//				traceMiddleware.Build(),
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user/:id/detail",
//			mdls: []Middleware{
//				traceMiddleware.Build(),
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user/:id",
//			mdls: []Middleware{
//				prometheusMiddleware.Build(),
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/order/detail",
//			mdls: []Middleware{
//				traceMiddleware.Build(),
//				prometheusMiddleware.Build(),
//			},
//		},
//		// 通配符测试用例
//		{
//			method: http.MethodGet,
//			path:   "/order/*",
//			mdls: []Middleware{
//				errhdlMiddleware,
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/order/*/info",
//			mdls: []Middleware{
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/order/detail/info",
//			mdls: []Middleware{
//				errhdlMiddleware,
//			},
//		},
//		{
//			method: http.MethodGet,
//			path:   "/order/:id/*",
//			mdls: []Middleware{
//			},
//		},
//	}
//
//	for _,v:=range testRoutes{
//		server.UseV1(v.method,v.path,v.mdls...)
//	}
//}
