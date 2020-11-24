package chat

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Server struct {
	id1 int
	id2 int
	id3 int
	w1  int
	w2  int
	w3  int
}

func (s *Server) SayHello(ctx context.Context, message *Response) (*Message, error) {
	// write to disk
	fileName := "lmhv_" + message.Info
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write/save buffer to disk
	ioutil.WriteFile(fileName, message.FileChunk, os.ModeAppend)

	fmt.Println("Split to : ", fileName)
	return &Message{Body: ""}, nil
}

func (s *Server) SayHello2(ctx context.Context, message *Message) (*Message, error) {
	// write to disk
	fmt.Println("Confirmas?")

	return &Message{In: 1}, nil
}
func (s *Server) GenerarPropuesta(ctx context.Context, message *Message) (*Propuesta, error) {
	cantidad := int(message.In)
	var propuesta Propuesta

	for i := 0; i < cantidad; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 3
		chosendn := rand.Intn(max-min+1) + min
		switch chosendn {
		case 1:
			propuesta.Id1++
		case 2:
			propuesta.Id2++
		case 3:
			propuesta.Id3++
		}
	}

	return &propuesta, nil
}
func (s *Server) PedirConfirmacion(ctx context.Context, message *Message) (*Message, error) {
	count := 0
	ganador := int(message.In)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := NewChatServiceClient(conn3)
	defer conn3.Close()

	var responde *Message
	message2 := Message{
		Body: "u2u",
	}
	if ganador == 1 {
		responde, _ = c2.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}
		responde, _ = c3.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}

	}
	if ganador == 2 {
		responde, _ = c.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}
		responde, _ = c3.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}

	}
	if ganador == 3 {
		responde, _ = c.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}
		responde, _ = c2.SayHello2(context.Background(), &message2)
		if responde.In == 1 {
			count++
		}

	}
	if count == 2 {
		return &Message{In: 1}, nil
	} else {
		return &Message{In: 0}, nil

	}
}

func (s *Server) Repartir(ctx context.Context, message *Message) (*Message, error) {
	ganador := int(message.In)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := NewChatServiceClient(conn3)
	defer conn3.Close()

	var responde *Message
	message2 := Message{
		Body: "u2u",
	}
	if ganador == 1 {
		responde, _ = c.SayHello2(context.Background(), &message2)
	}
	if ganador == 2 {
		responde, _ = c2.SayHello2(context.Background(), &message2)
	}
	if ganador == 3 {
		responde, _ = c3.SayHello2(context.Background(), &message2)
	}
	fmt.Println("Uwu : ", responde)

	return &Message{Body: "Holi"}, nil
}
