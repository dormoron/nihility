package nihility

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	// Start 启动服务器
	Start(add string) error

	AddRoute(method string, path string, handler HandleFunc)
}

type HTTPServer struct {
}

var _ Server = &HTTPServer{}

// ServeHTTP HTTpServer 处理请求入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Request: request,
		Writer:  writer,
	}
	s.serve(ctx)
}

func (s *HTTPServer) Start(add string) error {
	return http.ListenAndServe(add, s)
}

func (s *HTTPServer) GET(path string, handler HandleFunc) {
	s.AddRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) POST(path string, handler HandleFunc) {
	s.AddRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) PUT(path string, handler HandleFunc) {
	s.AddRoute(http.MethodPut, path, handler)
}

func (s *HTTPServer) DELETE(path string, handler HandleFunc) {
	s.AddRoute(http.MethodDelete, path, handler)
}

func (s *HTTPServer) OPTIONS(path string, handler HandleFunc) {
	s.AddRoute(http.MethodOptions, path, handler)
}

func (s *HTTPServer) AddRoute(method string, path string, handler HandleFunc) {

}

func (s *HTTPServer) serve(ctx *Context) {

}
