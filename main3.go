package main

import (
	"log"
	"net"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func con() {
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

func main() {
	go con()
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}

	defer conn2.Close()
	for {

	}
}
