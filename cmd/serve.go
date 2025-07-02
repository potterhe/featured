/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/potterhe/featured/internal/server"
	helloworldpb "github.com/potterhe/featured/proto/helloworld"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServerEndpoint = "localhost:50051"
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

		lis, err := net.Listen("tcp", grpcServerEndpoint)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		helloworldpb.RegisterGreeterServer(s, &server.Server{})
		reflection.Register(s)
		log.Printf("server listening at %v", lis.Addr())
		go func() {
			log.Fatalln(s.Serve(lis))
		}()

		gwmux := runtime.NewServeMux()
		// Register Greeter
		ctx := context.Background()
		otps := []grpc.DialOption{grpc.WithInsecure()}
		helloworldpb.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, grpcServerEndpoint, otps)
		if err != nil {
			log.Fatalln("Failed to register gateway:", err)
		}

		gwServer := &http.Server{
			Addr:    ":8090",
			Handler: gwmux,
		}

		log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
		log.Fatalln(gwServer.ListenAndServe())
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
