package nihility

import (
	"context"
	"github.com/nothingZero/nihility/registry"
	"google.golang.org/grpc/resolver"
	"time"
)

type grpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func InitGrpcResolver(r registry.Registry, timeout time.Duration) (resolver.Builder, error) {
	return &grpcResolverBuilder{
		r:       r,
		timeout: timeout,
	}, nil
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &grpcResolver{
		clientConn: cc,
		registry:   g.r,
		target:     target,
		timeout:    g.timeout,
	}
	r.resolve()
	go r.watch()
	return r, nil
}

func (g *grpcResolverBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	target     resolver.Target
	registry   registry.Registry
	clientConn resolver.ClientConn
	timeout    time.Duration
	close      chan struct{}
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instances, err := g.registry.ListServices(ctx, g.target.Endpoint())
	if err != nil {
		g.clientConn.ReportError(err)
		return
	}
	address := make([]resolver.Address, 0, len(instances))
	for _, instance := range instances {
		address = append(address, resolver.Address{Addr: instance.Address})
	}
	err = g.clientConn.UpdateState(resolver.State{
		Addresses: address,
	})

	if err != nil {
		g.clientConn.ReportError(err)
		return
	}

}

func (g *grpcResolver) watch() {
	events, err := g.registry.Subscribe(g.target.Endpoint())
	if err != nil {
		g.clientConn.ReportError(err)
		return
	}
	for {
		select {
		case <-events:
			g.resolve()
		case <-g.close:
			return
		}
	}
}

func (g *grpcResolver) Close() {
	close(g.close)
}
