package main

import (
	"context"
	"log"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}

	defer conn.Close()
	c := chat.NewChatServiceClient(conn)

	message := chat.Message{
		Body: "hola",
	}
	response, err := c.SayHello(context.Background(), &message)
	log.Printf("Holi 1 %s", response.Body)
}
