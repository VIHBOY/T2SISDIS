package main

import (
	"context"
	"log"
	"net"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func con() {
	var iniciar chat.Message
	iniciar.Body = "0"

	lis, err := net.Listen("tcp", "dist27:9002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := chat.Server{}
	s.CambiarRA(context.Background(), &iniciar)
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("a %v", err)
	}

}

func main() {
	go con()

	for {

	}
}
