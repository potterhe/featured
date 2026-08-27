package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/potterhe/featured/proto/helloworld"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var meter = otel.Meter("greeter")
var helloCnt metric.Int64Counter

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
}

func init() {
	initMetrics()
	helloCnt, _ = meter.Int64Counter("hello")
}

func initMetrics() error {
	// 2. 创建 Prometheus Exporter —— 关键：它会自动注册到默认的 /metrics 端点
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("create prometheus exporter: %w", err)
	}

	// 3. 用这个 exporter 构建 MeterProvider（全局）
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)
	return nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	helloCnt.Add(context.TODO(), 1)

	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
