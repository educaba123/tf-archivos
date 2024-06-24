package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var IPReg string
var IPSend string
var searchDni int

type Sale struct {
	Date             string  `json:"date"`
	Salesperson      string  `json:"salesperson"`
	CustomerName     string  `json:"customer_name"`
	Dni              int     `json:"dni"`
	CarMake          string  `json:"car_make"`
	CarModel         string  `json:"car_model"`
	CarYear          int     `json:"car_year"`
	SalePrice        float64 `json:"sale_price"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionEarned float64 `json:"commission_earned"`
}

var sales []Sale
var ch, ch2 chan (string)

func cargarDatos() {
	url := "https://raw.githubusercontent.com/educaba123/tf-archivos/main/dataset_sales.csv"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error haciendo la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error al descargar el archivo CSV: %s", resp.Status)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error leyendo el archivo CSV: %v", err)
	}

	for _, record := range records[1:] { // Saltar encabezado
		dni, _ := strconv.Atoi(record[3])
		carYear, _ := strconv.Atoi(record[6])
		salePrice, _ := strconv.ParseFloat(record[7], 64)
		commissionRate, _ := strconv.ParseFloat(record[8], 64)
		commissionEarned, _ := strconv.ParseFloat(record[9], 64)

		sale := Sale{
			Date:             record[0],
			Salesperson:      record[1],
			CustomerName:     record[2],
			Dni:              dni,
			CarMake:          record[4],
			CarModel:         record[5],
			CarYear:          carYear,
			SalePrice:        salePrice,
			CommissionRate:   commissionRate,
			CommissionEarned: commissionEarned,
		}
		sales = append(sales, sale)
	}
}

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

func resuelveListar(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(sales, "", " ")
	io.WriteString(res, string(jsonBytes))
	log.Println("Llamada al endpoint /list")
}

func resuelveBuscarCliente(res http.ResponseWriter, req *http.Request) {
	log.Println("Llamada al endpoint /sale")
	res.Header().Set("Content-Type", "application/json")
	dniStr := req.URL.Query().Get("dni")
	dni, err := strconv.Atoi(dniStr)
	if err != nil {
		http.Error(res, "Invalid DNI", http.StatusBadRequest)
		return
	}
	for _, sale := range sales {
		if sale.Dni == dni {
			jsonBytes, _ := json.MarshalIndent(sale, "", " ")
			io.WriteString(res, string(jsonBytes))
			return
		}
	}
	http.Error(res, "Venta no encontrada", http.StatusNotFound)
}

func recibirDni(res http.ResponseWriter, req *http.Request) {
	log.Println("Llamada al endpoint /search")
	if req.URL.Path != "/search" {
		http.Error(res, "404 not found.", http.StatusNotFound)
		return
	}

	switch req.Method {
	case "GET":
		io.WriteString(res, `
			<!doctype html>
			<html>
			<head><title>Search DNI</title></head>
			<body>
			<form action="/search" method="post">
				<label for="dni">DNI:</label>
				<input type="text" id="dni" name="dni">
				<input type="submit" value="Search">
			</form>
			</body>
			</html>
		`)
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() err: %v", err)
			return
		}
		dni := req.FormValue("dni")
		searchDni, _ = strconv.Atoi(dni)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{"message": fmt.Sprintf("DNI %d registrado para su búsqueda", searchDni)})
	default:
		fmt.Fprintf(res, "Sorry, only GET and POST methods are supported.")
	}
}

func añadirCliente(res http.ResponseWriter, req *http.Request) {
	log.Println("Llamada al endpoint /add")
	if req.URL.Path != "/add" {
		http.Error(res, "404 not found.", http.StatusNotFound)
		return
	}

	switch req.Method {
	case "GET":
		http.ServeFile(res, req, "form.html")
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() err: %v", err)
			return
		}
		fecha := req.FormValue("fecha")
		vendedor := req.FormValue("vendedor")
		cliente := req.FormValue("cliente")
		dni, _ := strconv.Atoi(req.FormValue("dni"))
		marca := req.FormValue("marca")
		modelo := req.FormValue("modelo")
		año, _ := strconv.Atoi(req.FormValue("año"))
		precio, _ := strconv.ParseFloat(req.FormValue("precio"), 64)
		comision_porc, _ := strconv.ParseFloat(req.FormValue("comision_porc"), 64)
		comision_gan, _ := strconv.ParseFloat(req.FormValue("comision_gan"), 64)

		nuevaVenta := Sale{
			Date:             fecha,
			Salesperson:      vendedor,
			CustomerName:     cliente,
			Dni:              dni,
			CarMake:          marca,
			CarModel:         modelo,
			CarYear:          año,
			SalePrice:        precio,
			CommissionRate:   comision_porc,
			CommissionEarned: comision_gan,
		}
		sales = append(sales, nuevaVenta)

		fmt.Fprintf(res, "Nuevo registro añadido con éxito\n")
		cliente1 := cliente
		enviarParametros(cliente1)
	default:
		fmt.Fprintf(res, "Sorry, only GET and POST methods are supported.")
	}
}

func enviarParametros(par string) string {
	ip := descubrirIP()
	remotehost := fmt.Sprintf("%s:%d", ip, 8000)
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()

	fmt.Fprintln(con, par)
	bf := bufio.NewReader(con)
	msg, _ := bf.ReadString('\n')
	return msg
}

func manejadorRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/list", resuelveListar)
	http.HandleFunc("/sale", resuelveBuscarCliente)
	http.HandleFunc("/search", recibirDni)
	http.HandleFunc("/add", añadirCliente)

	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese la dirección IP de registro: ")
	ip, _ := bufferIn.ReadString('\n')
	ip = strings.TrimSpace(ip)
	fmt.Print("Ingrese el puerto de registro: ")
	port, _ := bufferIn.ReadString('\n')
	port = strings.TrimSpace(port)
	IPReg = fmt.Sprintf("%s:%s", ip, port)

	log.Fatal(http.ListenAndServe(IPReg, nil))
}

func main() {
	ch = make(chan string)
	ch2 = make(chan string)

	cargarDatos()
	manejadorRequest()
}
