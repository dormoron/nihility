package nihility

import (
	"context"
	"github.com/nothingZero/nihility/registry"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServerOptions func(server *Server)

type Server struct {
	name            string
	registry        registry.Registry
	registryTimeout time.Duration
	*grpc.Server
	listener net.Listener
}

func InitServer(name string, opts ...ServerOptions) (*Server, error) {
	res := &Server{
		name:   name,
		Server: grpc.NewServer(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func WithRegister(r registry.Registry) ServerOptions {
	return func(server *Server) {
		server.registry = r
	}
}

// Start 服务已经准备好
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	// 有注册中心
	if s.registry != nil {
		// 超时时间
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		// 注册
		err = s.registry.Register(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: listener.Addr().String(),
		})
		if err != nil {
			return err
		}
	}
	err = s.Serve(listener)
	return err

}

func (s *Server) Close() error {
	if s.registry != nil {
		err := s.registry.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}
