package main

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/salirezam/grpc_client_server_demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
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

func startGRPCServer(address, certFile, keyFile string) error {
	var connections []*api.Connection
	// create a listnere on TCP port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create a server
	server := api.Server{connections}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// Create gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.StreamInterceptor(streamInterceptor)}

	// create gRPC server
	grpcServer := grpc.NewServer(opts...)

	// attach the Broadcast service to the server
	api.RegisterBroadcastServer(grpcServer, &server)

	// attach the Greeting service to the server
	api.RegisterGreetingServer(grpcServer, &server)

	// start the server
	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %s", err)
	}

	return nil
}

func startRESTServer(address, grpcAddress, certFile string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}
	// Setup the client gRPC options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	// Register greeting
	err = api.RegisterGreetingHandlerFromEndpoint(ctx, mux, grpcAddress, opts)

	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}
	log.Printf("starting HTTP/1.1 REST server on %s", address)
	http.ListenAndServe(address, mux)
	return nil
}

// run a gRPC
func main() {
	// server addresses
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	// certificate files
	certFile := "cert/server.crt"
	keyFile := "cert/server.key"

	// start gRPC server
	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// start REST server
	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// infinite loop
	log.Printf("Entering infinite loop")
	select {}
}
