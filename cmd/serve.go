/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/potterhe/featured/internal/server"
	pb "github.com/potterhe/featured/proto/helloworld"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Set up Prometheus exporter
		exporter, err := prometheus.New()
		if err != nil {
			log.Fatalf("failed to create prometheus exporter: %v", err)
		}
		provider := metric.NewMeterProvider(metric.WithReader(exporter))
		otel.SetMeterProvider(provider)
		defer func() {
			if err := provider.Shutdown(ctx); err != nil {
				log.Printf("failed to shutdown meter provider: %v", err)
			}
		}()

		// Set up stdout log exporter
		logExporter, err := stdoutlog.New()
		if err != nil {
			log.Fatalf("failed to create log exporter: %v", err)
		}
		lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)))
		global.SetLoggerProvider(lp)
		slog.SetDefault(otelslog.NewLogger("featured"))
		defer func() {
			if err := lp.Shutdown(ctx); err != nil {
				log.Printf("failed to shutdown logger provider: %v", err)
			}
		}()

		// Set up stdout trace exporter
		traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			log.Fatalf("failed to create trace exporter: %v", err)
		}
		tp := trace.NewTracerProvider(trace.WithBatcher(traceExporter))
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{}))
		defer func() {
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("failed to shutdown tracer provider: %v", err)
			}
		}()

		// Start Prometheus metrics HTTP server
		go func() {
			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.Handler())
			metricsAddr := ":9090"
			log.Printf("prometheus metrics server listening at %s/metrics", metricsAddr)
			if err := http.ListenAndServe(metricsAddr, mux); err != nil {
				log.Fatalf("failed to start metrics server: %v", err)
			}
		}()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
		pb.RegisterGreeterServer(s, &server.Server{})
		reflection.Register(s)

		// Graceful shutdown on signal
		go func() {
			<-ctx.Done()
			log.Println("shutting down gRPC server...")
			s.GracefulStop()
		}()

		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
