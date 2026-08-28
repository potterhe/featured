package server

import (
	"context"
	"log"

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
func (s *Server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	helloCnt.Add(context.TODO(), 1)

	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
