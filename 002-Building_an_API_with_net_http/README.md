# Building an API with the "net/http" Package
---
## The Net Package
![GoLang Net Package](/images/002_golang-net-package.png)
Package net provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets.    
[Net Package Documentation](https://pkg.go.dev/net)

#### Basic TCP Input/Output

The Dial function opens a connection to a server `net.Dial("tcp", "golang.org:80")`    
[Example: Dial Function](Example-Dial_Function.go)    

The Listen function creates server listeners to recieve connections `net.Listen("tcp", ":8000")`    
[Example: Listen Function](Example-Listen_Function.go)    

After a connection is opened between our client and server, input/output can be sent or recieved on both sides    

## net/http

While it is possible to use basic input/output to handle/send hand written http requests, there is no need to recreate the wheel, the [net/http](https://pkg.go.dev/net/http) package has everything we need for http input/output without needing to write all the protocols bits and pieces.    
[net/http Package Documentation](https://pkg.go.dev/net/http)    

The Get function creates a connection to the server and performs a get request `resp, err := http.Get("https://golang.org")`    
[Example: HTTP Get Function](Example-HTTP_Get_Function.go)    

The ListenAndServe function creates a server and listens for incoming connections, serving data from functions mapped to paths via the HandleFunc function    
```go
// register path handlers first
http.HandleFunc("/", myHandler)

// then start the server
log.Fatal(http.ListenAndServe(":8080", nil))
```

paths must be mapped to functions before starting the server because the ListenAndServe function blocks further execution   
[Example: HTTP Server](Example-HTTP_Server.go)    

#### Routing

**What is a Mux**    
A mux or a Multiplexer, is an http request router.    
when using functions like HandleFunc or ListenAndServ, if you dont apply a custom mux, then go uses the default global mux    
this is fine in simple applications however when you need to do more complex routing you will want to use a custom built mux    
there are security concerns around using the default global mux since it is shared globally, which means routes can potentially conflict with imported packages, a custom mux is separate from the global scope and has the routes specified to it.    

To create a custom mux we use `router := http.NewServeMux()`    
We can then call functions from the new mux like `router.HandleFunc()`    

**Request Methods**    
by default, if unspecified, all methods are valid in a route, however we may wanna limit them, or do different things with different methods.    
we can do this by specifying which method we want before the routes path like so `router.HandleFunc("GET /path", Handler)`    

**Dynamic Path Variables**    
sometimes we may want to use a dynamic path and pull information into a variable.   
we can do this by adding brackets to the part of the path we want as a variable like so
```go
// {id} is the dynamic variable
router.HandleFunc("GET /products/{id}", func(w http.ResponseWriter, r *http.Request) {
    // access the value using PathValue
    productID := r.PathValue("id")
})
```


#### Middleware

Middleware is software that is run while processing a request, it takes an `http.handler` and returns a new `http.handler`.    
Commonly used for injecting headers, validating authentication, and logging requests.    


Middleware can be nested and chained together like so
```go
/*
    Manual chaining (executes in order: Logger -> Auth -> Handler)
*/

// for specific paths
http.HandleFunc("/path", Logger(Auth(Handler)))

// for every request
http.ListenAndServe(":8080", Logger(Auth(Handler)))
```
[Example: Middleware](Example-Middleware.go)

#### JSON in Go
In Go, JSON handling is managed through the "encoding/json" package. The primary methods for interacting with JSON are Marshalling (converting Go types to JSON) and Unmarshalling (converting JSON to Go types)    
To convert a Go struct into JSON, use `json.Marshal`. To parse JSON back into a struct, use `json.Unmarshal`    
[Example: Marshalling and Unmarshalling JSON](Example-Marshalling.go)   


#### Putting it all Together

we now have all the necessary bits and pieces needed to create a REST API.    
Normally we would want to connect a database to handle our data but for simplicity we will store it in memory    
the completed project looks like:    

```go
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
```
[This Code](Finished_Project.go)
