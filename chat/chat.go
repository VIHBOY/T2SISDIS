package chat

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
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
	fmt.Println("Total : ", message.Cantidad)
	propuesta, _ := GenerarPropuesta(message.Cantidad)
	fmt.Println("Split to : ", propuesta.Id1)
	return &Message{Body: "Holi"}, nil
}

func (s *Server) SayHello2(ctx context.Context, message *Message) (*Message, error) {
	// write to disk

	return &Message{Body: "Holi"}, nil
}
func GenerarPropuesta(c uint64) (*Propuesta, error) {
	cantidad := int(c)
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
