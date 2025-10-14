package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
}

// http.Handler interface implementation
func (appServer *AppServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		// w.Write([]byte("Hello, World!\n"))
		fmt.Fprintln(w, "Hello, World!")
	case "/products":
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
		}
	case "/customers":
		fmt.Fprintln(w, "The list of customers will be served")
	default:
		http.Error(w, "resource not found", http.StatusNotFound)
	}

}

func main() {
	appServer := &AppServer{}
	if err := http.ListenAndServe(":8080", appServer); err != nil {
		log.Println(err)
	}
}
