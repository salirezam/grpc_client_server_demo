package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/salirezam/grpc_client_server_demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

var client api.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *api.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &api.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(str api.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}

			fmt.Printf("%v (%s): %s\n", msg.GetId(), msg.GetUser().GetName(), msg.GetMessage())
		}
	}(stream)

	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make(chan int)

	name := flag.String("name", "Alireza", "Name of the user")
	flag.Parse()

	id := sha256.Sum256([]byte(timestamp.String() + *name))

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect: %s", err)
	}
	//defer conn.Close()

	//client := api.NewGreetingClient(conn)
	client = api.NewBroadcastClient(conn)
	user := &api.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	connect(user)
	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &api.ChatMessage{
				Id:        user.Id,
				Message:   scanner.Text(),
				User:      user,
				Timestamp: timestamp.String(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error Sending Message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done

	// response, err := client.SayHello(context.Background(), &api.Message{Body: "Hey!"})
	// if err != nil {
	// 	log.Fatalf("Can't call SayHello: %s", err)
	// }
	// log.Printf("Response from server: %s", response.Body)

}
