package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

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

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}

	defer conn2.Close()
	c2 := chat.NewChatServiceClient(conn2)

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}

	defer conn3.Close()
	c3 := chat.NewChatServiceClient(conn3)

	fileToBeChunked := "./Los_Miserables-Hugo_Victor.pdf"

	file, err := os.Open(fileToBeChunked)

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 3
	chosendn := rand.Intn(max-min+1) + min
	fmt.Println(chosendn)

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
		case 2:
			response, _ = c2.SayHello(context.Background(), &message)
		case 3:
			response, _ = c3.SayHello(context.Background(), &message)

		}
		log.Printf("Holi 1 %s", response.Body)
	}

}
