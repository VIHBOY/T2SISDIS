package main

import (
	"log"
	"net"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := chat.Server{}
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("a %v", err)
	}

}
