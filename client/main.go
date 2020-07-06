package main

import (
	"github.com/salirezam/grpc_client_server_demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect: %s", err)
	}
	defer conn.Close()

	client := api.NewGreetingClient(conn)

	response, err := client.SayHello(context.Background(), &api.Message{Body: "Hey!"})
	if err != nil {
		log.Fatalf("Can't call SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)

}
