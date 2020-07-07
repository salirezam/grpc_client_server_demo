# Simple gRPC Client Server Chat

This GO chat application is based on gRPC using Protobuf. It has two part which are a client that sends and receives messages and a server which receives and broadcasts messages to all the clients. It supports SSL/TLS integration to secure the communication. Also, it has simple built-in authentication to identify the clients.

### Requirements
- **Protobuf :** You can find how to install it on Mac or Linux from the following [link](http://google.github.io/proto-lens/installing-protoc.html)
- **Go Lang :** You can find various installation tutorials for different platforms like the following [link](https://ahmadawais.com/install-go-lang-on-macos-with-homebrew/) for Mac
- You will also need to install Go dependencies

### Generating client and server code
- Run the following command to the gRPC client and server interfaces from our .proto service definition.<br/>
`protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/api.proto`
- Build the server with the following command.<br/>
`go build -i -v -o bin/server github.com/salirezam/grpc_client_server_demo/server`
- Build the client with the following command.<br/>
`go build -i -v -o bin/client github.com/salirezam/grpc_client_server_demo/client`

### Building certificates
Create a directory with name `cert` in the root directory and run the following commands to generate the required certificates.
Use `localhost` for host name or common name if you are using the default code.
```
$ openssl genrsa -out cert/server.key 2048
$ openssl req -new -x509 -sha256 -key cert/server.key -out cert/server.crt -days 3650
$ openssl req -new -sha256 -key cert/server.key -out cert/server.csr
$ openssl x509 -req -sha256 -in cert/server.csr -signkey cert/server.key -out cert/server.crt -days 3650
```
### Run client and server
- **Server :** `./bin/server`
- **Client :** `./bin/client -name [NAME] -username [USERNAME] -password [PASSWORD]`<br/>
There are two default username/password that can be used for testing purposes:<br/>
username:Alireza, password: 123456<br/>
username:John, password: 654321<br/>

### Useful Resources
- A basic tutorial introduction to gRPC in Go. [Link](https://grpc.io/docs/languages/go/basics/)
- How we use gRPC to build a client/server system in Go. [Link](https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2)