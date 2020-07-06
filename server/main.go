package main

import (
	"fmt"
	"github.com/salirezam/grpc_client_server_demo/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

// run a gRPC
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create a server
	server := api.Server{}

	// create gRPC server
	grpcServer := grpc.NewServer()

	// attach the Greeting service to the server
	api.RegisterGreetingServer(grpcServer, &server)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %s", err)
	}

}
