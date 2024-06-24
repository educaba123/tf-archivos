package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func descubrirIP() string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if strings.HasPrefix(i.Name, "Wi-Fi") || strings.HasPrefix(i.Name, "Ethernet") {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				switch t := addr.(type) {
				case *net.IPNet:
					if t.IP.To4() != nil {
						return t.IP.To4().String()
					}
				}
			}
		}
	}
	return "127.0.0.1"
}

func manejador(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	datos, _ := r.ReadString('\n')
	fmt.Printf("DATOS RECIBIDOS: %s\n", datos)
	time.Sleep(time.Millisecond * 1000)
	fmt.Fprintf(conn, "Msg = %s\n", "Dueño del vehículo: "+datos)
}

func server() {
	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del servidor: ")
	port, _ := bufferIn.ReadString('\n')
	port = strings.TrimSpace(port)
	ip := descubrirIP()
	address := fmt.Sprintf("%s:%s", ip, port)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		var a []any = []any{"Falla en la comunicación", err.Error()}
		fmt.Fprintln(os.Stdout, a...)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Servidor escuchando en", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Falla en la conexion", err.Error())
		}
		go manejador(conn)
	}
}

func main() {
	server()
}
