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

//Server is
/***
* struct Server
**
* Estructura de server
**
* Fields:
* sync.Mutex mu : Herramienta para sincronizacion
* int id1 : Id del DataNode1
* int id2 : Id del DataNode2
* int id3 : Id del DataNode3
* string ra : Estado del datanode con respecto a ocupar la zona critica
* []Response ListaChunks : Lista de chunks del datanode
***/
type Server struct {
	mu          sync.Mutex
	id1         int
	id2         int
	id3         int
	ra          string
	ListaChunks []Response
}

//contains is
/***
* func contains
**
* Verifica si en una lista ya existe un item en especfico
**
* Input:
* []int32 s : Lista cualquiera
* int32 e : Item a verificar
*
**
* Output:
* bool : retorna true o false, dependiendo de si existe o no el elemento
***/
func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//SayHello3 is
/***
* func SayHello3
**
* Recibe un chunk y lo escribe en un archivo
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Response message : Mensaje con estructura de tipo Response
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
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

//SayHello is
/***
* func SayHello
**
* Agrega un chunk a la lista de chunks
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Response message : Mensaje con estructura de tipo Response
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
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

//SayHello2 is
/***
* func SayHello2
**
* Imprime un confirmar y retorna un 1 en el cuerpo del request
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Response message : Mensaje con estructura de tipo Response
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) SayHello2(ctx context.Context, message *Message) (*Message, error) {
	// write to disk
	fmt.Println("Confirmas?")

	return &Message{In: 1}, nil
}

//GenerarPropuesta is
/***
* func GenerarPropuesta
**
* Genera una propuesta de distribucion de forma aleatoria
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Response message : Mensaje con estructura de tipo Response
*
**
* Output:
* *Propuesta : Mensaje con estructura de tipo Propuesta
* error : Constante nil
***/
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

//PedirConfirmacion is
/***
* func PedirConfirmacion
**
* Pide confirmar la propuesta a todos los datanodes
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Response message : Mensaje con estructura de tipo Response
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) PedirConfirmacion(ctx context.Context, message *Message) (*Message, error) {
	count := 0
	ganador := int(message.In)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		fmt.Println("uwu")

	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
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
	} else {
		message3.In = 0
	}
	return &message3, nil
}

//Repartir is
/***
* func Repartir
**
* Reparte los chunks a partir de la propuesta
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta propuesta : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) Repartir(ctx context.Context, propuesta *Propuesta) (*Message, error) {

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

//EscribirPropuesta is
/***
* func EscribirPropuesta
**
* Escribe la propuesta generada en el log del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta propuesta : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) EscribirPropuesta(ctx context.Context, propuesta *Propuesta) (*Message, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9004", grpc.WithInsecure())
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

//HelperEscribirPropuesta is
/***
* func HelperEscribirPropuesta
**
* Funcion de ayuda para escribir la propuesta generada en el log del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta propuesta : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
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

//EscribirPropuestaDis is
/***
* func EscribirPropuestaDis
**
* Escribe la propuesta generada en el log del NameNode,
* utilizando el algoritmo de Ricart y Agrawala
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta propuesta : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) EscribirPropuestaDis(ctx context.Context, propuesta *Propuesta) (*Message, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	fmt.Println(err)
	if err != nil {
		log.Fatalf("uwu %s", err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	fmt.Println(err2)

	if err2 != nil {
		log.Fatalf("uwu %s", err2)
	}
	c2 := NewChatServiceClient(conn2)

	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
	fmt.Println(err3)

	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c3 := NewChatServiceClient(conn3)

	defer conn3.Close()
	var auxiliar Message
	var respuesta1 *Message
	var respuesta2 *Message
	var respuesta3 *Message
	Nodisponible2 := []int32{}
	Sidisponible2 := []int32{}

	message2 := Message{
		Body: "u2u",
	}
	var flag int
	r1, _ := c.SayHello2(context.Background(), &message2)
	if r1 == nil {
		Nodisponible2 = append(Nodisponible2, 1)
	} else {
		respuesta1, _ = c.ConsultarRA(context.Background(), &auxiliar)

	}
	r, _ := c2.SayHello2(context.Background(), &message2)
	if r == nil {
		Nodisponible2 = append(Nodisponible2, 2)
	} else {
		Sidisponible2 = append(Sidisponible2, 2)
		respuesta2, _ = c2.ConsultarRA(context.Background(), &auxiliar)
		flag = 1
	}
	r2, _ := c3.SayHello2(context.Background(), &message2)
	if r2 == nil {
		Nodisponible2 = append(Nodisponible2, 3)
	} else {
		Sidisponible2 = append(Sidisponible2, 3)
		respuesta3, _ = c3.ConsultarRA(context.Background(), &auxiliar)
		if flag == 1 {
			flag = 2
		} else {
			flag = 0
		}
	}
	if len(Sidisponible2) == 2 && flag == 2 {
		for {
			if respuesta1.Body == "0" && respuesta2.Body == "0" && respuesta3.Body == "0" {
				fmt.Printf("XD\n")
				auxiliar.Body = "1"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "2":
					c2.CambiarRA(context.Background(), &auxiliar)
				case "3":
					c3.CambiarRA(context.Background(), &auxiliar)
				}
				var conn4 *grpc.ClientConn
				conn4, err4 := grpc.Dial(":9004", grpc.WithInsecure())
				if err4 != nil {
					log.Fatalf("uwu %s", err)
				}
				c4 := NewChatServiceClient(conn4)

				defer conn.Close()

				s.mu.Lock()
				c4.HelperEscribirPropuesta(context.Background(), propuesta)
				s.mu.Unlock()

				auxiliar.Body = "0"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "2":
					c2.CambiarRA(context.Background(), &auxiliar)
				case "3":
					c3.CambiarRA(context.Background(), &auxiliar)
				}
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}
	if len(Sidisponible2) == 1 && flag == 1 {
		for {
			if respuesta1.Body == "0" && respuesta2.Body == "0" {
				fmt.Printf("XD\n")
				auxiliar.Body = "1"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "2":
					c2.CambiarRA(context.Background(), &auxiliar)
				}
				var conn4 *grpc.ClientConn
				conn4, err4 := grpc.Dial(":9004", grpc.WithInsecure())
				if err4 != nil {
					log.Fatalf("uwu %s", err)
				}
				c4 := NewChatServiceClient(conn4)

				defer conn.Close()

				s.mu.Lock()
				c4.HelperEscribirPropuesta(context.Background(), propuesta)
				s.mu.Unlock()

				auxiliar.Body = "0"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "2":
					c2.CambiarRA(context.Background(), &auxiliar)
				}
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}
	if len(Sidisponible2) == 1 && flag == 0 {
		for {
			if respuesta1.Body == "0" && respuesta3.Body == "0" {
				fmt.Printf("XD\n")
				auxiliar.Body = "1"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "3":
					c3.CambiarRA(context.Background(), &auxiliar)
				}
				var conn4 *grpc.ClientConn
				conn4, err4 := grpc.Dial(":9004", grpc.WithInsecure())
				if err4 != nil {
					log.Fatalf("uwu %s", err)
				}
				c4 := NewChatServiceClient(conn4)

				defer conn.Close()

				s.mu.Lock()
				c4.HelperEscribirPropuesta(context.Background(), propuesta)
				s.mu.Unlock()

				auxiliar.Body = "0"
				switch propuesta.Chosendn {
				case "1":
					c.CambiarRA(context.Background(), &auxiliar)
				case "3":
					c3.CambiarRA(context.Background(), &auxiliar)
				}
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}

	return &Message{Body: ""}, nil
}

//HelperEscribirPropuestaDis is
/***
* func HelperEscribirPropuestaDis
**
* Funcion de ayuda para escribir la propuesta generada en el log del NameNode
* en base al algoritmo de Ricart y Agrawala
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta propuesta : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) HelperEscribirPropuestaDis(ctx context.Context, propuesta *Propuesta) (*Message, error) {
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

//AgregarTitulo is
/***
* func AgregarTitulo
**
* Agrega libro al log de titulos del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) AgregarTitulo(ctx context.Context, message *Message) (*Message, error) {

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9004", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c4 := NewChatServiceClient(conn3)
	defer conn3.Close()
	c4.HelperAgregarTitulo(context.Background(), message)
	return &Message{Body: ""}, nil
}

//HelperAgregarTitulo is
/***
* func HelperAgregarTitulo
**
* Funcion de ayuda para agregar libros al log de titulos del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
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

//VerTitulos is
/***
* func VerTitulos
**
* Muestra todos los libros del log de titulos del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Titulos : Mensaje con estructura de tipo Titulos
* error : Constante nil
***/
func (s *Server) VerTitulos(ctx context.Context, message *Message) (*Titulos, error) {

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9004", grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("uwu %s", err3)
	}
	c4 := NewChatServiceClient(conn3)
	defer conn3.Close()
	Titulos, _ := c4.HelperVerTitulos(context.Background(), message)
	return Titulos, nil
}

//HelperVerTitulos is
/***
* func HelperVerTitulos
**
* Funcion de ayuda para mostrar todos los libros del log de titulos del NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Titulos : Mensaje con estructura de tipo Titulos
* error : Constante nil
***/
func (s *Server) HelperVerTitulos(ctx context.Context, message *Message) (*Titulos, error) {
	var Titulos Titulos
	file, err := os.Open("Titulos.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Titulos.Titulos = append(Titulos.Titulos, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &Titulos, nil
}

//ObtenerUbicaciones is
/***
* func ObtenerUbicaciones
**
* Obtiene las ubicaciones de los chunks
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Titulos : Mensaje con estructura de tipo Titulos
* error : Constante nil
***/
func (s *Server) ObtenerUbicaciones(ctx context.Context, message *Message) (*Titulos, error) {
	tituloabuscar := message.Body
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
		if a[0] == tituloabuscar {
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

//BuscarChunks is
/***
* func BuscarChunks
**
* Busca los chunks en los DataNodes
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Titulos ti : Mensaje con estructura de tipo Titulos
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) BuscarChunks(ctx context.Context, ti *Titulos) (*Message, error) {
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
	s.mu.Lock()
	for _, b := range ti.Titulos {
		a := strings.Split(b, " ")
		fmt.Println(a[0])
		i, _ := strconv.Atoi(a[1])
		switch i {

		case 2:
			fmt.Println("2")
			chunk, _ := c2.HacerChunks(context.Background(), &Message{Body: a[0]})
			if chunk == nil {
				fmt.Println("No se Pudo Obtener " + a[0])
			} else {
				c.SayHello3(context.Background(), chunk)
			}

		case 3:
			fmt.Println("3")
			chunk, _ := c3.HacerChunks(context.Background(), &Message{Body: a[0]})
			if chunk == nil {
				fmt.Println("No se Pudo Obtener " + a[0])
			} else {
				c.SayHello3(context.Background(), chunk)
			}
		}

	}

	s.mu.Unlock()
	return &Message{Body: ""}, nil

}

//HacerChunks is
/***
* func HacerChunks
**
* Crea chunks en base a una estructura de tipo Message
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message ti : Mensaje con estructura de tipo Message
*
**
* Output:
* *Response : Mensaje con estructura de tipo Response
* error : Constante nil
***/
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

//Unir is
/***
* func Unir
**
* Busca chunks y los une en un archivo
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message ti : Mensaje con estructura de tipo Message
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
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

		fmt.Println("Bip Bup Trabajando en  : [", writePosition, "] bytes Bip Bup *Sonidos de robot*")
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

		fmt.Println("LLevamos ", n, " bytes Asombroso :o")

		fmt.Println("La parte N°[", j, "] esta siendo añadida al archivo: ", newFileName+" Velocito")
	}

	// now, we close the newFileName
	file.Close()
	return &Message{Body: ""}, nil
}

//GenerarPropuesta2 is
/***
* func GenerarPropuesta2
**
* Genera una segunda propuesta de distribucion
* sin considerar los nodos caidos
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Propuesta : Mensaje con estructura de tipo Propuesta
* error : Constante nil
***/
func (s *Server) GenerarPropuesta2(ctx context.Context, message *Message) (*Propuesta, error) {
	cantidad := int(message.In)
	fmt.Println(cantidad)
	a := message.Nodisponible
	var propuesta Propuesta
	if cantidad == 3 {
		if len(a) == 2 {
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
		}
		if len(a) == 1 {
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
		}
		if len(a) == 0 {
			propuesta.Id1++
			propuesta.L1 = append(propuesta.L1, int32(0))
			propuesta.Pos = append(propuesta.Pos, 1)
			propuesta.Id2++
			propuesta.L2 = append(propuesta.L2, int32(1))
			propuesta.Pos = append(propuesta.Pos, 2)
			propuesta.Id3++
			propuesta.L3 = append(propuesta.L3, int32(2))
			propuesta.Pos = append(propuesta.Pos, 3)
		}
		return &propuesta, nil

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
			}
			if len(a) == 1 {
				if chosendn == int(a[0]) {
					chosendn = rand.Intn(max-min+1) + min
				}
				if chosendn != int(a[0]) {
					break
				}
			}
			if len(a) == 0 {
				break
			}

		}
		fmt.Println(chosendn)
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

//PedirConfirmacion2 is
/***
* func PedirConfirmacion2
**
* Pide confirmar distribución sin considerar los nodos caidos
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta message : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) PedirConfirmacion2(ctx context.Context, message *Propuesta) (*Message, error) {
	/*a := message.Nodisponible
	count := 0
	ganador := int(message.In)
	*/
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	c := NewChatServiceClient(conn)

	defer conn.Close()

	var conn2 *grpc.ClientConn
	conn2, err2 := grpc.Dial(":9001", grpc.WithInsecure())
	if err2 != nil {
		fmt.Println("uwu")

	}
	c2 := NewChatServiceClient(conn2)
	defer conn2.Close()

	var conn3 *grpc.ClientConn
	conn3, err3 := grpc.Dial(":9002", grpc.WithInsecure())
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
	for _, b := range message.Pos {
		switch b {
		case 1:
			responde, _ = c.SayHello2(context.Background(), &message2)
			if responde == nil {
				message3.In = 0
				b := contains(message3.Nodisponible, int32(1))
				if b == false {
					message3.Nodisponible = append(message3.Nodisponible, int32(1))
				}

			}
		case 2:
			responde, _ = c2.SayHello2(context.Background(), &message2)
			if responde == nil {
				message3.In = 0
				b := contains(message3.Nodisponible, int32(2))
				if b == false {
					message3.Nodisponible = append(message3.Nodisponible, int32(2))
				}
			}
		case 3:
			responde, _ = c3.SayHello2(context.Background(), &message2)
			if responde == nil {
				message3.In = 0
				b := contains(message3.Nodisponible, int32(3))
				if b == false {
					message3.Nodisponible = append(message3.Nodisponible, int32(3))
				}
			}
		}

	}
	if len(message3.Nodisponible) > 0 {
		message3.In = 0
		return &message3, nil
	}
	message3.In = 1
	return &message3, nil
	/*if ganador == 1 {
		if len(a) == 2 {
			if a[0] == 2 && a[1] == 3 {
				message3.In = 1
				return &message3, nil
			}
		}
		if len(a) == 1 {
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
		if len(a) == 0 {
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

	}
	if ganador == 2 {
		if len(a) == 2 {
			if a[0] == 1 && a[1] == 3 {
				message3.In = 1
				return &message3, nil
			}
		}
		if len(a) == 1 {
			if a[0] != 1 {
				fmt.Println("uwu1")
				responde, _ = c.SayHello2(context.Background(), &message2)
				if responde != nil && responde.In == 1 {
					count++
				} else {
					message3.Nodisponible = append(message3.Nodisponible, int32(1))
				}
			}
			if a[0] != 3 {
				fmt.Println("uwu12")

				responde, _ = c3.SayHello2(context.Background(), &message2)
				if responde != nil && responde.In == 1 {
					count++
				} else {
					message3.Nodisponible = append(message3.Nodisponible, int32(3))
				}
			}
		}
		if len(a) == 0 {
			responde, _ = c.SayHello2(context.Background(), &message2)
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

	}
	if ganador == 3 {
		if len(a) == 2 {
			if a[0] == 2 && a[1] == 3 {
				message3.In = 1
				return &message3, nil
			}
		}
		if len(a) == 1 {
			if a[0] != 2 {
				responde, _ = c2.SayHello2(context.Background(), &message2)
				if responde != nil && responde.In == 1 {
					count++
				} else {
					message3.Nodisponible = append(message3.Nodisponible, int32(2))
				}
			}
			if a[0] != 1 {
				responde, _ = c.SayHello2(context.Background(), &message2)
				if responde != nil && responde.In == 1 {
					count++
				} else {
					message3.Nodisponible = append(message3.Nodisponible, int32(1))
				}
			}
		}
		if len(a) == 0 {
			responde, _ = c2.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(2))
			}
			responde, _ = c.SayHello2(context.Background(), &message2)
			if responde != nil && responde.In == 1 {
				count++
			} else {
				message3.Nodisponible = append(message3.Nodisponible, int32(1))
			}
		}

	}
	fmt.Println(count)
	if count == 2-len(a) {
		message3.In = 1
		return &message3, nil
	} else {
		message3.In = 0
		return &message3, nil

	}*/
}

//PedirConfirmacionNM is
/***
* func PedirConfirmacionNM
**
* Pide confirmar distribución al NameNode
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Propuesta message : Mensaje con estructura de tipo Propuesta
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) PedirConfirmacionNM(ctx context.Context, message *Propuesta) (*Message, error) {
	var conn4 *grpc.ClientConn
	conn4, err := grpc.Dial(":9004", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	c4 := NewChatServiceClient(conn4)
	defer conn4.Close()

	message2 := Message{
		Body: "u2u",
	}
	r, _ := c4.SayHello2(context.Background(), &message2)
	if r.In == 1 {
		return &Message{In: 1}, nil
	}
	return &Message{In: 0}, nil

}

//CambiarRA is
/***
* func CambiarRA
**
* Cambia el estado del DataNode, especificando si este
* se encuentra en una zona critica o no
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) CambiarRA(ctx context.Context, message *Message) (*Message, error) {
	s.ra = message.Body
	fmt.Println(message.Body)
	return &Message{Body: message.Body}, nil
}

//ConsultarRA is
/***
* func ConsultarRA
**
* Consulta el estado del DataNode, especificando si este
* se encuentra en una zona critica o no
**
* Input:
* context.Context ctx : Contexto de alguna coneccion
* *Message message : Mensaje con estructura de tipo Message
*
**
* Output:
* *Message : Mensaje con estructura de tipo Message
* error : Constante nil
***/
func (s *Server) ConsultarRA(ctx context.Context, message *Message) (*Message, error) {
	return &Message{Body: s.ra}, nil
}
