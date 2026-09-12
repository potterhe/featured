package server

import (
	"context"
	"log/slog"

	pb "github.com/potterhe/featured/proto/helloworld"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("greeter")
var helloCnt metric.Int64Counter

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
}

func init() {
	helloCnt, _ = meter.Int64Counter("hello")
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	helloCnt.Add(ctx, 1)

	slog.InfoContext(ctx, "received request", slog.String("name", in.GetName()))
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
