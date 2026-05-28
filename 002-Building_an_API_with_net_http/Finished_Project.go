package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"time"
)

// Product structure
type Product struct {
	Id  string  `json:"Id"`
	Title   string  `json:"Title"`
	Desc	string  `json:"desc"`
	Price   int64   `json:"price"`
}

// slice to store Products
var Products []Product

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Products)
}

func returnSingleProduct(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")

	for _, product := range Products {
		if product.Id == key {
			json.NewEncoder(w).Encode(product)
		}
	}
}


func createNewProduct(w http.ResponseWriter, r *http.Request) {
	// convert post body to a product
	reqBody, _ := io.ReadAll(r.Body)
	var product Product
	json.Unmarshal(reqBody, &product)

	// add product to products
	Products = append(Products, product)

	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	for index, product := range Products {
		if product.Id == id {
			Products = append(Products[:index], Products[index+1:]...)
		}
	}

}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {

	Products = []Product{
		Product{Id: "1", Title: "Product 1", Desc: "Product Description", Price: 100},
		Product{Id: "2", Title: "Product 2", Desc: "Product Description", Price: 200},
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /products", returnAllProducts)
	router.HandleFunc("POST /product", createNewProduct)
	router.HandleFunc("DELETE /product/{id}", deleteProduct)
	router.HandleFunc("GET /product/{id}", returnSingleProduct)

	log.Println("Launching Server on Port 3000")
	log.Fatal(http.ListenAndServe(":3000", logger(router)))

}
