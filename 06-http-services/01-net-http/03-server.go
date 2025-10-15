package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sdrapkin/guid"
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

type HandlerFunction func(http.ResponseWriter, *http.Request)
type MiddlwareFunction func(HandlerFunction) HandlerFunction

type AppServer struct {
	routes      map[string]HandlerFunction
	middlewares []MiddlwareFunction
}

func NewAppServer() *AppServer {
	return &AppServer{
		routes: make(map[string]HandlerFunction),
	}
}

func (appServer *AppServer) AddRoute(path string, handlerFn HandlerFunction) {
	for i := len(appServer.middlewares) - 1; i >= 0; i-- {
		handlerFn = appServer.middlewares[i](handlerFn)
	}
	appServer.routes[path] = handlerFn
}

func (appServer *AppServer) UseMiddleware(middlewareFn MiddlwareFunction) {
	appServer.middlewares = append(appServer.middlewares, middlewareFn)
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

func logMiddleware(handlerFn HandlerFunction) HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			log.Printf("%x - %s - %s - %v\n", r.Context().Value("req-id"), r.Method, r.URL.Path, elapsed)
		}()
		handlerFn(w, r)
	}
}

func timeoutMiddleware(handlerFn HandlerFunction) HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		timeoutCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		handlerFn(w, r.WithContext(timeoutCtx))
		if timeoutCtx.Err() != nil {
			log.Printf("%x - %s\n", r.Context().Value("req-id"), "request timed out")
		}
	}
}

func requestIdMiddleware(handlerFn HandlerFunction) HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		idCtx := context.WithValue(r.Context(), "req-id", guid.New())
		handlerFn(w, r.WithContext(idCtx))
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second) // simlating a time consuming operation
	select {
	case <-r.Context().Done():
		http.Error(w, "request timed out", http.StatusRequestTimeout)
	default:
		fmt.Fprintln(w, "Hello, World!")
	}
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
	appServer.UseMiddleware(requestIdMiddleware)
	appServer.UseMiddleware(timeoutMiddleware)
	appServer.UseMiddleware(logMiddleware)
	appServer.AddRoute("/", IndexHandler)
	appServer.AddRoute("/products", ProductsHandler)
	appServer.AddRoute("/customers", CustomersHandler)
	if err := http.ListenAndServe(":8080", appServer); err != nil {
		log.Println(err)
	}
}
