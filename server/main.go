package main

import (
	"fmt"
	"github.com/salirezam/grpc_client_server_demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

// run a gRPC
func main() {
	var connections []*api.Connection

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", 7777))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create a server
	server := api.Server{connections}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// Create gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// create gRPC server
	grpcServer := grpc.NewServer(opts...)

	// attach the Greeting service to the server
	api.RegisterGreetingServer(grpcServer, &server)

	// attach the Broadcast service to the server
	api.RegisterBroadcastServer(grpcServer, &server)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %s", err)
	}

}
