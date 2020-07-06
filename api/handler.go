package api

import (
	"golang.org/x/net/context"
	"log"
)

// gRPC server representation
type Server struct{}

// SayHello generates response to a Greeting request
func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive greeting message %s", in.Body)
	return &Message{Body: "Hello! Hola!"}, nil
}

func (s *Server) mustEmbedUnimplementedGreetingServer() {}
