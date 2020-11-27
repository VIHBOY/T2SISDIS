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
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	fmt.Println(text2)
	Agregar2 := chat.Message{
		Body: text3,
	}
	Ti, _ := c.ObtenerUbicaciones(context.Background(), &Agregar2)
	fmt.Println(Ti)
	c.BuscarChunks(context.Background(), Ti)
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
	fmt.Println(propuesta.Pos)
	fmt.Println(propuesta.L1)
	fmt.Println(propuesta.L2)
	fmt.Println(propuesta.L3)

	if respuesta.In == 1 {
		fmt.Println("Propuesta Aceptada")
	}

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
