package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VIHBOY/T2SISDIS/chat"
	"google.golang.org/grpc"
)

func con() {
	lis, err := net.Listen("tcp", "dist25:9000")
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
func do2(text string, c chat.ChatServiceClient, c2 chat.ChatServiceClient, c3 chat.ChatServiceClient) {
	Agregar := chat.Message{
		Body: "text",
	}
	Titulos, _ := c.VerTitulos(context.Background(), &Agregar)
	fmt.Printf("Titulos Disponibles\n")
	for _, a := range Titulos.Titulos {
		fmt.Printf("Titulos: %s\n", a)
	}
	fmt.Println("Ingrese Titulo")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese Archivo a bajar")
	text2, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text3 := strings.Replace(text, "\n", "", -1)
	fmt.Println(text3)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist28:9004", grpc.WithInsecure())
	c4 := chat.NewChatServiceClient(conn)

	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	fmt.Println(text2)
	Agregar2 := chat.Message{
		Body: text3,
	}
	Ti, _ := c4.ObtenerUbicaciones(context.Background(), &Agregar2)
	fmt.Println(Ti)
	c.BuscarChunks(context.Background(), Ti)
	me := chat.Message{
		Body: text3,
		In:   int32(len(Ti.Titulos)),
	}
	c.Unir(context.Background(), &me)
	defer conn.Close()

}
func do(text string, c chat.ChatServiceClient, c2 chat.ChatServiceClient, c3 chat.ChatServiceClient) {

	/*var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)

	defer conn3.Close()*/

	fileToBeChunked := "./" + text + ".pdf"
	titulo := text
	file, err := os.Open(fileToBeChunked)

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 3
	chosendn := 1
	/*rand.Intn(max-min) + min*/
	message2 := chat.Message{
		Body: "u2u",
	}
	for {
		switch chosendn {
		case 2:
			r, _ := c2.SayHello2(context.Background(), &message2)
			fmt.Println(r)
			if r == nil {
				chosendn = rand.Intn(max-min) + min
				fmt.Println("POTO")
				fmt.Println(chosendn)

			} else {
				break
			}
		case 3:
			r, _ := c3.SayHello2(context.Background(), &message2)
			if r == nil {
				chosendn = rand.Intn(max-min) + min
			} else {
				break
			}
		}
		fmt.Println("POTO2")
		fmt.Println(chosendn)
		break
	}

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
	propuesta.Titulo = titulo

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
	fmt.Println("uwu")
	fmt.Println(respuesta.Nodisponible)
	fmt.Println(propuesta.Pos)
	fmt.Println(propuesta.L1)
	fmt.Println(propuesta.L2)
	fmt.Println(propuesta.L3)

	_ = respuesta
	/*
		if respuesta.In == 1 {
			fmt.Println("Propuesta Aceptada")
		}
	*/

	for {
		if respuesta.In == 1 {
			fmt.Println("Propuesta Aceptada")
			break
		} else {
			fmt.Println("Propuesta No Aceptada")
			fmt.Println(respuesta.Nodisponible)
			Message_totalchunks := chat.Message{
				In:           int32(totalPartsNum),
				Nodisponible: respuesta.Nodisponible,
			}
			switch chosendn {
			case 1:
				propuesta, _ = c.GenerarPropuesta2(context.Background(), &Message_totalchunks)
			case 2:
				propuesta, _ = c2.GenerarPropuesta2(context.Background(), &Message_totalchunks)
			case 3:
				propuesta, _ = c3.GenerarPropuesta2(context.Background(), &Message_totalchunks)
			}
			Aviso := chat.Message{
				In:           int32(chosendn),
				Nodisponible: respuesta.Nodisponible,
			}
			var respuesta *chat.Message
			switch chosendn {
			case 1:
				respuesta, _ = c.PedirConfirmacion2(context.Background(), &Aviso)
			case 2:
				respuesta, _ = c2.PedirConfirmacion2(context.Background(), &Aviso)
			case 3:
				respuesta, _ = c3.PedirConfirmacion2(context.Background(), &Aviso)

			}
			fmt.Println(respuesta.In)
			if respuesta.In == 1 {
				break
			}
		}
	}
	propuesta.Titulo = titulo

	fmt.Println(propuesta.Pos)
	fmt.Println(propuesta.L1)
	fmt.Println(propuesta.L2)
	fmt.Println(propuesta.L3)

	switch chosendn {
	case 1:
		respuesta, _ = c.EscribirPropuesta(context.Background(), propuesta)
	case 2:
		respuesta, _ = c2.EscribirPropuesta(context.Background(), propuesta)
	case 3:
		respuesta, _ = c3.EscribirPropuesta(context.Background(), propuesta)

	}

	fmt.Printf("Dividir en  %d Chunks.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		message := chat.Response{
			Info:      strconv.FormatUint(i, 10),
			Name:      titulo,
			Elegido:   1,
			Cantidad:  totalPartsNum,
			FileChunk: partBuffer,
		}

		switch chosendn {
		case 1:
			c.SayHello(context.Background(), &message)
		case 2:
			c2.SayHello(context.Background(), &message)
		case 3:
			c3.SayHello(context.Background(), &message)
		}
	}
	Agregar := chat.Message{
		Body: text,
	}
	var Titulos *chat.Titulos
	switch chosendn {
	case 1:
		c.Repartir(context.Background(), propuesta)
		c.AgregarTitulo(context.Background(), &Agregar)
		Titulos, _ = c.VerTitulos(context.Background(), &Agregar)
	case 2:
		c2.Repartir(context.Background(), propuesta)
		c2.AgregarTitulo(context.Background(), &Agregar)
		Titulos, _ = c2.VerTitulos(context.Background(), &Agregar)

	case 3:
		c3.Repartir(context.Background(), propuesta)
		c3.AgregarTitulo(context.Background(), &Agregar)
		Titulos, _ = c3.VerTitulos(context.Background(), &Agregar)
	}
	fmt.Println(Titulos)
	fmt.Printf("Titulos Disponibles\n")
	for _, a := range Titulos.Titulos {
		fmt.Printf("Titulos: %s\n", a)
	}
}

func main() {
	go con()
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	fmt.Println(err)
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := chat.NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	fmt.Println(err2)

	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := chat.NewChatServiceClient(conn2)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	fmt.Println(err3)

	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)

	defer conn3.Close()
	for {
		fmt.Println("-------------------------------------------")
		fmt.Println("Menú")
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Ingrese Archivo a subir")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		fmt.Println(text)
		do(text, c, c2, c3)
		do2(text, c, c2, c3)

	}

}
