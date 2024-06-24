package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	dataURL = "https://raw.githubusercontent.com/educaba123/tf-archivos/main/dataset_sales.csv"
)

type Point struct {
	AssessedValue float64
	SaleAmount    float64
}

type Result struct {
	Node int
	M    float64
	B    float64
}

func leerDatosDesdeURL(url string) ([]Point, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al descargar el archivo CSV: %s", resp.Status)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var points []Point
	for _, record := range records[1:] { // Saltar encabezado
		assessedValue, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, err
		}
		saleAmount, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, Point{AssessedValue: assessedValue, SaleAmount: saleAmount})
	}
	return points, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingresa la dirección del nodo: ")
	direccion, _ := reader.ReadString('\n')
	direccion = strings.TrimSpace(direccion)

	fmt.Print("Con cuántos nodos se va a unir: ")
	var n int
	fmt.Scanf("%d\n", &n)

	addrs := make([]string, n)
	for i := range addrs {
		fmt.Printf("Nodo %d: ", i+1)
		addrs[i], _ = reader.ReadString('\n')
		addrs[i] = strings.TrimSpace(addrs[i])
	}

	points, err := leerDatosDesdeURL(dataURL)
	if err != nil {
		fmt.Println("Error al leer los datos desde el archivo:", err)
		return
	}

	chunkSize := len(points) / n
	results := make(chan Result, n)

	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == n-1 {
			end = len(points)
		}

		go func(node int, addr string, points []Point) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println("Error al conectar con el nodo:", err)
				return
			}
			defer conn.Close()

			for _, p := range points {
				fmt.Fprintf(conn, "%f %f\n", p.AssessedValue, p.SaleAmount)
			}
			fmt.Fprintln(conn, "EOF")

			var m, b float64
			fmt.Fscanf(conn, "%f %f\n", &m, &b)
			results <- Result{Node: node, M: m, B: b}
		}(i+1, addrs[i], points[start:end])
	}

	var allResults []string
	for i := 0; i < n; i++ {
		result := <-results
		fmt.Printf("Resultados del nodo %d: m = %f, b = %f\n", result.Node, result.M, result.B)
		allResults = append(allResults, fmt.Sprintf("Nodo %d: m = %f, b = %f\n", result.Node, result.M, result.B))
	}

	fmt.Print("Ingresa la dirección del servidor (IP:PUERTO): ")
	serverAddr, _ := reader.ReadString('\n')
	serverAddr = strings.TrimSpace(serverAddr)

	sendResultsToServer(serverAddr, strings.Join(allResults, ""))
}

func sendResultsToServer(serverAddr, data string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error al conectar con el servidor:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", data)
	fmt.Println("Resultados enviados al servidor")
}
