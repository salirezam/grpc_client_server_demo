package api

import (
	"golang.org/x/net/context"
	glog "google.golang.org/grpc/grpclog"
	"log"
	"os"
	"sync"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

// gRPC Server representation
type Server struct {
	Connection []*Connection
}

// Connection struct
type Connection struct {
	stream Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

// SayHello generates response to a Greeting request
func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive greeting message %s", in.Body)
	return &Message{Body: "Hello! Hola!"}, nil
}

// CreateStream initiates connection to Broadcast Server and return a stream of message
func (s *Server) CreateStream(pconn *Connect, stream Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

// BroadcastMessage sends message to all user
func (s *Server) BroadcastMessage(ctx context.Context, msg *ChatMessage) (*Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *ChatMessage, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Info("Sending message to: ", conn.stream)

				if err != nil {
					grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &Close{}, nil
}

func (s *Server) mustEmbedUnimplementedGreetingServer()  {}
func (s *Server) mustEmbedUnimplementedBroadcastServer() {}
