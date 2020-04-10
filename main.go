package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"transactiond/txmsg"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SendMessage(ctx context.Context, in *txmsg.RpcMessage) (*txmsg.RpcMessage, error) {
	log.Printf("Received: %v", in.GroupId)
	return &txmsg.RpcMessage{}, nil
}

func main() {

	port := "127.0.0.1:6102"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	txmsg.RegisterManageServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
