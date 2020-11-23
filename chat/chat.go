package chat

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
)

type Server struct {
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
	return &Message{Body: "Holi"}, nil
}
