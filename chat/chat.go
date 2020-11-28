package chat

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	sync "sync"
	"time"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Server struct {
	mu          sync.Mutex
	id1         int
	id2         int
	id3         int
	w1          int
	w2          int
	w3          int
	ListaChunks []Response
}

func (s *Server) SayHello3(ctx context.Context, message *Response) (*Message, error) {

	// write to disk
	fileName := message.Name + "_" + message.Info
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write/save buffer to disk
	ioutil.WriteFile(fileName, message.FileChunk, os.ModeAppend)

	fmt.Println("Recibiste: ", fileName)

	return &Message{Body: ""}, nil
}
func (s *Server) SayHello(ctx context.Context, message *Response) (*Message, error) {

	// write to disk

	s.ListaChunks = append(s.ListaChunks, Response{
		Info:      message.GetInfo(),
		Name:      message.GetName(),
		Elegido:   message.GetElegido(),
		Cantidad:  message.GetCantidad(),
		FileChunk: message.GetFileChunk(),
	})

	// write/save buffer to disk

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
	if cantidad == 3 {
		propuesta.Id1++
		propuesta.L1 = append(propuesta.L1, int32(0))
		propuesta.Pos = append(propuesta.Pos, 1)
		propuesta.Id2++
		propuesta.L2 = append(propuesta.L2, int32(1))
		propuesta.Pos = append(propuesta.Pos, 2)
		propuesta.Id3++
		propuesta.L3 = append(propuesta.L3, int32(2))
		propuesta.Pos = append(propuesta.Pos, 3)
		return &propuesta, nil
	}
	for i := 0; i < cantidad; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 3
		chosendn := rand.Intn(max-min+1) + min
		switch chosendn {
		case 1:
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(i))
			propuesta.Pos = append(propuesta.Pos, 1)
		case 2:
			propuesta.Id2++
			propuesta.L2 = append(propuesta.L2, int32(i))
			propuesta.Pos = append(propuesta.Pos, 2)

		case 3:
			propuesta.Id3++
			propuesta.L3 = append(propuesta.L3, int32(i))
			propuesta.Pos = append(propuesta.Pos, 3)

		}
	}

	return &propuesta, nil
}
func (s *Server) PedirConfirmacion(ctx context.Context, message *Message) (*Message, error) {
	count := 0
	ganador := int(message.In)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	if err2 != nil {
		fmt.Println("uwu")

	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	if err3 != nil {
		fmt.Println(err)
	}
	c3 := NewChatServiceClient(conn3)
	defer conn3.Close()

	var responde *Message
	message2 := Message{
		Body: "u2u",
	}
	message3 := Message{}
	if ganador == 1 {
		responde, _ = c2.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(2))
		}
		responde, _ = c3.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(3))
		}

	}
	if ganador == 2 {
		responde, _ = c.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(1))
		}
		responde, _ = c3.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(3))
		}

	}
	if ganador == 3 {
		responde, _ = c.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(1))
		}
		responde, _ = c2.SayHello2(context.Background(), &message2)
		if responde != nil && responde.In == 1 {
			count++
		} else {
			message3.Nodisponible = append(message3.Nodisponible, int32(2))
		}

	}

	if count == 2 {
		message3.In = 1
		return &message3, nil
	} else {
		message3.In = 0
		return &message3, nil

	}
}

func (s *Server) Repartir(ctx context.Context, propuesta *Propuesta) (*Message, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := NewChatServiceClient(conn2)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := NewChatServiceClient(conn3)

	defer conn3.Close()
	s.mu.Lock()

	for _, i := range propuesta.L1 {
		c.SayHello3(context.Background(), &s.ListaChunks[i])

	}
	for _, i := range propuesta.L2 {
		c2.SayHello3(context.Background(), &s.ListaChunks[i])

	}
	for _, i := range propuesta.L3 {
		c3.SayHello3(context.Background(), &s.ListaChunks[i])

	}
	s.ListaChunks = make([]Response, 0)

	s.mu.Unlock()
	return &Message{Body: ""}, nil
}
func (s *Server) EscribirPropuesta(ctx context.Context, propuesta *Propuesta) (*Message, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist28:9004", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	s.mu.Lock()
	c.HelperEscribirPropuesta(context.Background(), propuesta)
	s.mu.Unlock()
	return &Message{Body: ""}, nil
}
func (s *Server) HelperEscribirPropuesta(ctx context.Context, propuesta *Propuesta) (*Message, error) {
	f, err := os.OpenFile("Log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	totalchunks := len(propuesta.L1) + len(propuesta.L2) + len(propuesta.L3)
	encabezado := propuesta.Titulo + " " + strconv.Itoa(totalchunks)
	if _, err := f.WriteString(encabezado + "\n"); err != nil {
		log.Println(err)
	}
	for i2, i := range propuesta.Pos {
		if _, err := f.WriteString(propuesta.Titulo + "_" + strconv.Itoa(i2) + " " + strconv.Itoa(int(i)) + "\n"); err != nil {
			log.Println(err)
		}
	}
	return &Message{Body: ""}, nil
}
func (s *Server) AgregarTitulo(ctx context.Context, message *Message) (*Message, error) {

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist28:9004", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c4 := NewChatServiceClient(conn3)
	defer conn3.Close()
	c4.HelperAgregarTitulo(context.Background(), message)
	return &Message{Body: ""}, nil
}
func (s *Server) HelperAgregarTitulo(ctx context.Context, message *Message) (*Message, error) {
	f, err := os.OpenFile("Titulos.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(message.Body + "\n"); err != nil {
		log.Println(err)
	}
	return &Message{Body: ""}, nil
}
func (s *Server) VerTitulos(ctx context.Context, message *Message) (*Titulos, error) {

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist28:9004", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c4 := NewChatServiceClient(conn3)
	defer conn3.Close()
	Titulos, _ := c4.HelperVerTitulos(context.Background(), message)
	return Titulos, nil
}
func (s *Server) HelperVerTitulos(ctx context.Context, message *Message) (*Titulos, error) {
	var Titulos Titulos
	file, err := os.Open("Titulos.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		Titulos.Titulos = append(Titulos.Titulos, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &Titulos, nil
}
func (s *Server) ObtenerUbicaciones(ctx context.Context, message *Message) (*Titulos, error) {
	titulo_a_buscar := message.Body
	var Titulos Titulos
	file, err := os.Open("Log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lines []string
	var i2 int
	var linea int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for c, b := range lines {
		a := strings.Split(b, " ")
		if a[0] == titulo_a_buscar {
			linea = c
			i2, _ = strconv.Atoi(a[1])

		}
	}

	for i := 0; i < i2; i++ {
		Titulos.Titulos = append(Titulos.Titulos, lines[linea+i+1])
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &Titulos, nil

}
func (s *Server) BuscarChunks(ctx context.Context, ti *Titulos) (*Message, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := NewChatServiceClient(conn2)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := NewChatServiceClient(conn3)

	defer conn3.Close()
	s.mu.Lock()
	for _, b := range ti.Titulos {
		a := strings.Split(b, " ")
		fmt.Println(a[0])
		i, _ := strconv.Atoi(a[1])
		switch i {

		case 2:
			fmt.Println("2")
			chunk, _ := c2.HacerChunks(context.Background(), &Message{Body: a[0]})
			c.SayHello3(context.Background(), chunk)

		case 3:
			fmt.Println("3")
			chunk, _ := c3.HacerChunks(context.Background(), &Message{Body: a[0]})
			c.SayHello3(context.Background(), chunk)

		}

	}

	s.mu.Unlock()
	return &Message{Body: ""}, nil

}
func (s *Server) HacerChunks(ctx context.Context, ti *Message) (*Response, error) {
	file, err := os.Open(ti.Body)
	a := strings.Split(ti.Body, "_")
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	chunkInfo, err := file.Stat()
	var chunkSize int64 = chunkInfo.Size()
	chunkBufferBytes := make([]byte, chunkSize)
	reader := bufio.NewReader(file)
	_, err = reader.Read(chunkBufferBytes)

	message := Response{
		Name:      a[0],
		Info:      a[1],
		FileChunk: chunkBufferBytes,
	}
	fmt.Println(chunkBufferBytes)
	return &message, nil
}
func (s *Server) Unir(ctx context.Context, ti *Message) (*Message, error) {
	totalPartsNum := uint64(ti.In)

	// just for fun, let's recombine back the chunked files in a new file

	newFileName := ti.Body + "1.pdf"
	_, err := os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int64 = 0

	for j := uint64(0); j < totalPartsNum; j++ {

		//read a chunk
		currentChunkFileName := ti.Body + "_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

		n, err := file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		fmt.Println("Written ", n, " bytes")

		fmt.Println("Recombining part [", j, "] into : ", newFileName)
	}

	// now, we close the newFileName
	file.Close()
	return &Message{Body: ""}, nil
}
func (s *Server) GenerarPropuesta2(ctx context.Context, message *Message) (*Propuesta, error) {
	cantidad := int(message.In)
	a := message.Nodisponible
	var propuesta Propuesta
	if cantidad == 3 {
		if a[0] == 2 && a[1] == 3 {
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(0))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(1))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(2))
			propuesta.Pos = append(propuesta.Pos, 1)
			return &propuesta, nil

		}
		if a[0] == 2 {
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(0))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id3++
			propuesta.L3 = append(propuesta.L3, int32(1))
			propuesta.Pos = append(propuesta.Pos, 3)
			propuesta.Id3++
			propuesta.L3 = append(propuesta.L3, int32(2))
			propuesta.Pos = append(propuesta.Pos, 3)
		}
		if a[0] == 3 {
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(0))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(1))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id2++
			propuesta.L2 = append(propuesta.L2, int32(2))
			propuesta.Pos = append(propuesta.Pos, 2)
		}
		/*propuesta.Id1++
		propuesta.L1 = append(propuesta.L1, int32(0))
		propuesta.Pos = append(propuesta.Pos, 1)
		propuesta.Id2++
		propuesta.L2 = append(propuesta.L2, int32(1))
		propuesta.Pos = append(propuesta.Pos, 2)
		propuesta.Id3++
		propuesta.L3 = append(propuesta.L3, int32(2))
		propuesta.Pos = append(propuesta.Pos, 3)*/
		return &propuesta, nil
	}
	for i := 0; i < cantidad; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 3
		chosendn := rand.Intn(max-min+1) + min
		for {
			if len(a) == 2 {
				if chosendn == int(a[0]) || chosendn == int(a[1]) {
					chosendn = rand.Intn(max-min+1) + min
				}
				if chosendn != int(a[0]) && chosendn != int(a[1]) {
					break
				}
			} else {
				if chosendn == int(a[0]) {
					chosendn = rand.Intn(max-min+1) + min
				}
				if chosendn != int(a[0]) {
					break
				}
			}

		}
		switch chosendn {
		case 1:
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(i))
			propuesta.Pos = append(propuesta.Pos, 1)
		case 2:
			propuesta.Id2++
			propuesta.L2 = append(propuesta.L2, int32(i))
			propuesta.Pos = append(propuesta.Pos, 2)

		case 3:
			propuesta.Id3++
			propuesta.L3 = append(propuesta.L3, int32(i))
			propuesta.Pos = append(propuesta.Pos, 3)

		}
	}

	return &propuesta, nil
}
func (s *Server) PedirConfirmacion2(ctx context.Context, message *Message) (*Message, error) {
	a := message.Nodisponible
	count := 0
	ganador := int(message.In)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist25:9000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial("dist26:9001", grpc.WithInsecure())
	if err2 != nil {
		fmt.Println("uwu")

	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial("dist27:9002", grpc.WithInsecure())
	if err3 != nil {
		fmt.Println(err)
	}
	c3 := NewChatServiceClient(conn3)
	defer conn3.Close()

	var responde *Message
	message2 := Message{
		Body: "u2u",
	}
	message3 := Message{}
	if ganador == 1 {
		if len(a) == 2 {
			if a[0] == 2 && a[1] == 3 {
				message3.In = 1
				return &message3, nil
			}
		}
		if a[0] != 2 {
			responde, _ = c2.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(2))
			}
		}
		if a[0] != 3 {
			responde, _ = c3.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(3))
			}
		}

	}
	if ganador == 2 {
		if a[0] != 1 {
			responde, _ = c.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(1))
			}
		}
		if a[0] != 3 {
			responde, _ = c3.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(3))
			}
		}

	}
	if ganador == 3 {
		if a[0] != 1 {
			responde, _ = c.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(1))
			}
		}
		if a[0] != 2 {
			responde, _ = c2.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(2))
			}
		}

	}

	if count == 2-len(a) {
		message3.In = 1
		return &message3, nil
	} else {
		message3.In = 0
		return &message3, nil

	}
}
