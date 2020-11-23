package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func main() {

	fileToBeChunked := "./Los_Miserables-Hugo_Victor.pdf"

	file, err := os.Open(fileToBeChunked)

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

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("uwu %s", err)
		}

		defer conn.Close()
		c := chat.NewChatServiceClient(conn)

		message := chat.Response{
			Info:      strconv.FormatUint(i, 10),
			FileChunk: partBuffer,
		}
		response, err := c.SayHello(context.Background(), &message)
		log.Printf("Holi 1 %s", response.Body)
	}

}
