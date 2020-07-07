package main

import (
	"fmt"
	"github.com/salirezam/grpc_client_server_demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"strings"
)

// authenticateAgent check the clients login info
func authenticateClient(ctx context.Context) error {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		if (clientLogin == "Alireza" && clientPassword == "123456") ||
			(clientLogin == "John" && clientPassword == "654321") {
			log.Printf("authenticated client: %s", clientLogin)
		} else {
			log.Printf("bad username or password: %s", clientLogin)
			return fmt.Errorf("bad username or password")
		}

		return nil
	}

	return fmt.Errorf("missing credentials")
}

// streamInterceptor calls authenticateClient with current context
func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authenticateClient(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

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
	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.StreamInterceptor(streamInterceptor)}

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
