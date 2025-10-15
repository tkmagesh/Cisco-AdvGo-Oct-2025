package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Category string  `json:"category"`
}

var products []Product = []Product{
	{Id: 101, Name: "Pen", Cost: 10, Category: "Stationary"},
	{Id: 102, Name: "Pencil", Cost: 5, Category: "Stationary"},
	{Id: 103, Name: "Marker", Cost: 50, Category: "Stationary"},
}

type AppServer struct {
	routes map[string]func(http.ResponseWriter, *http.Request)
}

func NewAppServer() *AppServer {
	return &AppServer{
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (appServer *AppServer) AddRoute(path string, handlerFn func(http.ResponseWriter, *http.Request)) {
	appServer.routes[path] = handlerFn
}

// http.Handler interface implementation
func (appServer *AppServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handlerFn := appServer.routes[r.URL.Path]; handlerFn != nil {
		handlerFn(w, r)
		return
	}
	http.Error(w, "resource not found", http.StatusNotFound)
}

// application specific

// Modify the function to log the time taken for processing the request as well
func logMiddleware(handlerFn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s\n", r.Method, r.URL.Path)
		handlerFn(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "Hello, World!")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(products); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		var newProduct Product
		if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
			http.Error(w, "error parsing the request payload", http.StatusBadRequest)
			return
		}
		newProduct.Id = 101 + len(products)
		products = append(products, newProduct)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newProduct); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
	}
}

func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The list of customers will be served")
}

func main() {
	appServer := NewAppServer()
	appServer.AddRoute("/", logMiddleware(IndexHandler))
	appServer.AddRoute("/products", logMiddleware(ProductsHandler))
	appServer.AddRoute("/customers", logMiddleware(CustomersHandler))
	if err := http.ListenAndServe(":8080", appServer); err != nil {
		log.Println(err)
	}
}
