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
func do3(text string, c chat.ChatServiceClient, c2 chat.ChatServiceClient, c3 chat.ChatServiceClient) {

	/*var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)

	defer conn3.Close()*/
	Nodisponible2 := []int32{}
	Sidisponible2 := []int32{}

	message2 := chat.Message{
		Body: "u2u",
	}
	r1, _ := c.SayHello2(context.Background(), &message2)
	fmt.Println(r1)
	if r1 == nil {
		Nodisponible2 = append(Nodisponible2, 1)
	} else {
		Sidisponible2 = append(Sidisponible2, 1)
	}
	r, _ := c2.SayHello2(context.Background(), &message2)
	fmt.Println(r)
	if r == nil {
		Nodisponible2 = append(Nodisponible2, 2)
	} else {
		Sidisponible2 = append(Sidisponible2, 2)
	}
	r2, _ := c3.SayHello2(context.Background(), &message2)
	if r2 == nil {
		Nodisponible2 = append(Nodisponible2, 3)
	} else {
		Sidisponible2 = append(Sidisponible2, 3)
	}

	fileToBeChunked := "./" + text + ".pdf"
	titulo := text
	file, err := os.Open(fileToBeChunked)
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(Sidisponible2)
	chosendn := Sidisponible2[rand.Intn(max-min)+min]
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
	fmt.Println(chosendn)
	Message_totalchunks := chat.Message{
		In:           int32(totalPartsNum),
		Nodisponible: Nodisponible2,
	}
	var propuesta *chat.Propuesta
	switch chosendn {
	case 1:
		propuesta, _ = c.GenerarPropuesta2(context.Background(), &Message_totalchunks)
	case 2:
		propuesta, _ = c2.GenerarPropuesta2(context.Background(), &Message_totalchunks)
	case 3:
		propuesta, _ = c3.GenerarPropuesta2(context.Background(), &Message_totalchunks)
	}
	propuesta.Titulo = titulo
	propuesta.In = int32(chosendn)
	propuesta.Nodisponible = Nodisponible2

	var respuesta *chat.Message
	switch chosendn {
	case 1:
		respuesta, _ = c.PedirConfirmacionNM(context.Background(), propuesta)
	case 2:
		respuesta, _ = c2.PedirConfirmacionNM(context.Background(), propuesta)
	case 3:
		respuesta, _ = c3.PedirConfirmacionNM(context.Background(), propuesta)

	}
	fmt.Println("uwu")
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

			propuesta.Titulo = titulo
			propuesta.In = int32(chosendn)
			propuesta.Nodisponible = respuesta.Nodisponible
			var respuesta *chat.Message
			switch chosendn {
			case 1:
				respuesta, _ = c.PedirConfirmacionNM(context.Background(), propuesta)
			case 2:
				respuesta, _ = c2.PedirConfirmacionNM(context.Background(), propuesta)
			case 3:
				respuesta, _ = c3.PedirConfirmacionNM(context.Background(), propuesta)

			}
			fmt.Println(respuesta.In)
			if respuesta.In == 1 {
				break
			}
		}
	}

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

	Agregar := chat.Message{
		Body: text,
	}
	fmt.Println("Holi2")
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
	text3 := strings.Replace(text2, "\n", "", -1)
	fmt.Println(text3)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9004", grpc.WithInsecure())
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
	Nodisponible2 := []int32{}
	Sidisponible2 := []int32{}

	message2 := chat.Message{
		Body: "u2u",
	}
	r1, _ := c.SayHello2(context.Background(), &message2)
	if r1 == nil {
		Nodisponible2 = append(Nodisponible2, 1)
	} else {
		Sidisponible2 = append(Sidisponible2, 1)
	}
	r, _ := c2.SayHello2(context.Background(), &message2)
	if r == nil {
		Nodisponible2 = append(Nodisponible2, 2)
	} else {
		Sidisponible2 = append(Sidisponible2, 2)
	}
	r2, _ := c3.SayHello2(context.Background(), &message2)
	if r2 == nil {
		Nodisponible2 = append(Nodisponible2, 3)
	} else {
		Sidisponible2 = append(Sidisponible2, 3)
	}

	fileToBeChunked := "./" + text + ".pdf"
	titulo := text
	file, err := os.Open(fileToBeChunked)
	rand.Seed(time.Now().UnixNano())
	/*min := 0
	max := len(Sidisponible2)*/
	chosendn := 2
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
	fmt.Println("")

	fmt.Printf("Dividiendo el archivo en %d Chunks. Agradecemos su espera\n", totalPartsNum)

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
	fmt.Println("")

	Message_totalchunks := chat.Message{
		In:           int32(totalPartsNum),
		Nodisponible: Nodisponible2,
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
	fmt.Println(propuesta.Pos)
	fmt.Println(Nodisponible2)
	propuesta.Titulo = titulo
	propuesta.In = int32(chosendn)
	propuesta.Nodisponible = Nodisponible2

	var respuesta *chat.Message
	switch chosendn {
	case 1:
		respuesta, _ = c.PedirConfirmacion2(context.Background(), propuesta)
	case 2:
		respuesta, _ = c2.PedirConfirmacion2(context.Background(), propuesta)
	case 3:
		respuesta, _ = c3.PedirConfirmacion2(context.Background(), propuesta)

	}
	/*fmt.Println(propuesta.Pos)
	fmt.Println(propuesta.L1)
	fmt.Println(propuesta.L2)
	fmt.Println(propuesta.L3)*/

	_ = respuesta
	/*
		if respuesta.In == 1 {
			fmt.Println("Propuesta Aceptada")
		}
	*/

	for {
		if respuesta.In == 1 {
			fmt.Println("")
			fmt.Println("Propuesta Aceptada")
			fmt.Println("")

			break
		} else {
			fmt.Println("")

			fmt.Println("Propuesta No Aceptada")
			fmt.Println("")
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

			propuesta.Titulo = titulo
			propuesta.In = int32(chosendn)
			propuesta.Nodisponible = respuesta.Nodisponible
			var respuesta *chat.Message
			switch chosendn {
			case 1:
				respuesta, _ = c.PedirConfirmacion2(context.Background(), propuesta)
			case 2:
				respuesta, _ = c2.PedirConfirmacion2(context.Background(), propuesta)
			case 3:
				respuesta, _ = c3.PedirConfirmacion2(context.Background(), propuesta)

			}
			if respuesta.In == 1 {
				break
			}
		}
	}
	fmt.Println(propuesta.Pos)
	fmt.Println(Nodisponible2)
	fmt.Println("")
	fmt.Println("Escribiendo Propuesta en Log.txt Ubicado en otra maquina (magia!:D)")
	fmt.Println("")

	switch chosendn {
	case 1:
		respuesta, _ = c.EscribirPropuestaDis(context.Background(), propuesta)
	case 2:
		respuesta, _ = c2.EscribirPropuestaDis(context.Background(), propuesta)
	case 3:
		respuesta, _ = c3.EscribirPropuestaDis(context.Background(), propuesta)

	}

	Agregar := chat.Message{
		Body: text,
	}
	fmt.Println("")

	fmt.Println("El nodo aleatorio Asigando comenzara a repartir uwu")
	fmt.Println("")

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
	fmt.Println("")

	fmt.Printf("Titulos Disponibles\n")
	for _, a := range Titulos.Titulos {
		fmt.Printf("Titulos: %s\n", a)
	}
}

func main() {
	go con()
	var iniciar chat.Message
	iniciar.Body = "0"

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	fmt.Println(err)
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := chat.NewChatServiceClient(conn)
	c.CambiarRA(context.Background(), &iniciar)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	fmt.Println(err2)

	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := chat.NewChatServiceClient(conn2)
	c2.CambiarRA(context.Background(), &iniciar)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	fmt.Println(err3)

	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := chat.NewChatServiceClient(conn3)
	c3.CambiarRA(context.Background(), &iniciar)

	defer conn3.Close()
	for {
		fmt.Println("-------------------------------------------")
		fmt.Println("Menú")
		fmt.Println("")
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Ingrese Archivo a subir")
		fmt.Println("")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		fmt.Println("")
		fmt.Println("Subiendo:" + text + ".pdf")
		do(text, c, c2, c3)
		do2(text, c, c2, c3)

	}

}
