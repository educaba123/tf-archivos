package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Point struct {
	AssessedValue float64
	SaleAmount    float64
}

func linearRegression(points []Point) (m, b float64) {
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(points))
	for _, p := range points {
		sumX += p.AssessedValue
		sumY += p.SaleAmount
		sumXY += p.AssessedValue * p.SaleAmount
		sumX2 += p.AssessedValue * p.AssessedValue
	}
	if denom := (n*sumX2 - sumX*sumX); denom != 0 {
		m = (n*sumXY - sumX*sumY) / denom
		b = (sumY - m*sumX) / n
	}
	return
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var points []Point
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOF" {
			break
		}
		fields := strings.Fields(line)
		assessedValue, _ := strconv.ParseFloat(fields[0], 64)
		saleAmount, _ := strconv.ParseFloat(fields[1], 64)
		points = append(points, Point{AssessedValue: assessedValue, SaleAmount: saleAmount})
	}

	m, b := linearRegression(points)
	fmt.Fprintf(conn, "%f %f\n", m, b)
}

func main() {
	fmt.Print("Ingresa la dirección del nodo trabajador (ip:puerto): ")
	var direccion string
	fmt.Scanf("%s\n", &direccion)

	listener, err := net.Listen("tcp", direccion)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Nodo trabajador escuchando en", direccion)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar la conexión:", err)
			continue
		}
		go handleConnection(conn)
	}
}
