package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func con() {
	lis, err := net.Listen("tcp", ":9000")
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
	c := chat.NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := chat.NewChatServiceClient(conn2)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)

	defer conn3.Close()

	/*var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)

	defer conn3.Close()*/

	fileToBeChunked := "./Los_Miserables-Hugo_Victor.pdf"

	file, err := os.Open(fileToBeChunked)

	rand.Seed(time.Now().UnixNano())
	max := 3
	min := 1
	chosendn := rand.Intn(max-min) + min
	/*rand.Intn(max-min) + min*/

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 256000

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	Message_totalchunks := chat.Message{
		In: int32(totalPartsNum),
	}
	var propuesta *chat.Propuesta
	switch chosendn {
	case 1:
		propuesta, _ = c.GenerarPropuesta(context.Background(), &Message_totalchunks)
	case 2:
		propuesta, _ = c2.GenerarPropuesta(context.Background(), &Message_totalchunks)
	case 3:
		propuesta, _ = c3.GenerarPropuesta(context.Background(), &Message_totalchunks)

	}

	fmt.Println(propuesta)

	Aviso := chat.Message{
		In: int32(chosendn),
	}
	var respuesta *chat.Message
	switch chosendn {
	case 1:
		respuesta, _ = c.PedirConfirmacion(context.Background(), &Aviso)
	case 2:
		respuesta, _ = c2.PedirConfirmacion(context.Background(), &Aviso)
	case 3:
		respuesta, _ = c3.PedirConfirmacion(context.Background(), &Aviso)

	}
	if respuesta.In == 1 {
		fmt.Println("Propuesta Aceptada")

	}

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		message := chat.Response{
			Info:      strconv.FormatUint(i, 10),
			Name:      "Los_Miserables-Hugo_Victor.pdf",
			Elegido:   1,
			Cantidad:  totalPartsNum,
			FileChunk: partBuffer,
		}

		var response *chat.Message
		switch chosendn {
		case 1:
			response, _ = c.SayHello(context.Background(), &message)
			log.Printf("Aca %s", response.Body)

		case 2:
			response, _ = c2.SayHello(context.Background(), &message)
		case 3:
			response, _ = c3.SayHello(context.Background(), &message)
		}
	}

	for {

	}
}
