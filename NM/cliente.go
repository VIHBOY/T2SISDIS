package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func con() {
	lis, err := net.Listen("tcp", ":9004")
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
	fileName := "Log.txt"
	_, err4 := os.Create(fileName)

	if err4 != nil {
		fmt.Println(err4)
		os.Exit(1)
	}
	fileName2 := "Titulos.txt"
	_, err5 := os.Create(fileName2)

	if err5 != nil {
		fmt.Println(err5)
		os.Exit(1)
	}
	go con()
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}

	defer conn3.Close()
	for {

	}
}
