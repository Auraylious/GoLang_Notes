# Scalable API Routing
![Scalable API Routing](/images/004-scalable_apis.png)    
lets face it, putting all your routes in one file can be hectic to sort through, especially if you have other components in the same file like your middlewares, and your websocket logic etc etc.    

lets organize things for better scalability early on so we can better focus on the specific thing at hand when updating it later.   

### Why its bad to put everything into a single file
1. It takes longer to sift through all that extra code to get to where you wanna be
2. When testing new features, you have to make sure all the other stuff is in tact, usually by making a backup of the whole project or similar
3. When adding new routes, aesthetically you have to organize them into the code in a place where it makes sense, or live with messy code (shame).

### The Solution
were going to do it plugin style, where you drop in a file and it loads automatically without needing to change the core project files.    
our main file with our servers logic, middlewares and other components, will stay the same no matter how we change the routes.     

adding, disabling, and removing routes will be as simple as renaming/moving their specific file.    

### The Setup
our route files will make use of the go init() function to register themselves to a shared route index as soon as the package is imported, the index will later be iterated and applied to our router object (ServeMux)     

our project will be organized like this    
```
ModuleName/
    go.mod
    main.go
    api/
        main-api.go
        route-userapi.go
        route-whatever_you_want.go
```


the main function will look something like this (rewrite with your database implementation)    

- we define our imports and import the api sub package of our module    
- create a basic server, then we call the api Setup function instead of directly using `router.HandleFunc()`    
- usually we want a database to work with so ive included that logic as a placeholder    

```go
//main.go
package main

import (
    "log"
    "net/http"
    "some/database"
    "ModuleName/api"
)

func main() {

    // a basic server with a database connection

    db, _ := database.Open("logic to open a database connection")
    defer db.Close()

    router := http.NewServeMux()

    api.Setup(db, router) // we pass the db along with the router to the api, so that we can access the database from individual route files

    log.fatal(http.ListenAndServe(":3000", router))

}
```

then in our apis main.go file we will setup routines and variables for our route files to use, letting them register themselves as routes to be attached when Setup is called.    
after we set this up, this file remains relatively unchanged for the rest of development, letting us focus only on the specific routes were interested in.   


```go
//api/main-api.go
package api

import (
    "net/http"      // for type definitions
    "some/database" // for type definitions
)

// Where the Api is located
var ApiPath = "/api"

// RouteFactory: accepts dependencies and returns an HTTP handler function
type RouteHandler func(db *database.DB) http.HandlerFunc

// Define the Route Archetype
type Route struct {
        Method  string
        Pattern string
        Handler RouteHandler
}

// Place to Store Route Definitions
var Registry []Route

// adds a route to the global registry
func Register(method, pattern string, handler RouteHandler) {
        Registry = append(Registry, Route{
                Method:  method,
                Pattern: pattern,
                Handler: handler,
        })
}

// Binds each route to a Multiplexer
func Setup(db *database.DB, router *http.ServeMux) {
        for _, route := range Registry {
                handler := route.Handler(db)
                router.HandleFunc(route.Method+" "+route.Pattern, handler)
        }
}
```

Now the way go imports packages, all the Uppercase variables and definitions are available to the other files, so we can access the Register function from our other files.    

in this case the naming of our files wont effect the functionality, everything ending with .go gets initialized together when imported     
so for pure readability and organization we are naming these files like so "route-whatever.go"    


```
//api/route-whatever.go
package api

import (
    "fmt"
    "net/http"
    "some/database"
)

func init() { // always gets called while importing the api package, has access to global things such as the Register function

    // We Register Our Routes Here like so

    // GET
    Register(http.MethodGet, ApiPath + "/{id}", func (db *database.DB) http.HandlerFunc {
        return func (w http.ResponseWriter, r *http.Request){
            
            // this is where we do things when the route is called
            fmt.Fprintf(w, "Hello World " + r.PathValue("id") + "\n")
        }
    })

    // PUT
    Register(http.MethodPut, ApiPath + "/{id}", func (db *database.DB) http.HandlerFunc {
        return func (w http.ResponseWriter, r *http.Request){
            fmt.Fprintf(w, "Hello World " + r.PathValue("id") + "\n")
        }
    })

    // POST
    Register(http.MethodPost, ApiPath, func (db *database.DB) http.HandlerFunc {
        return func (w http.ResponseWriter, r *http.Request){
            fmt.Fprintf(w, "Hello World\n")
        }
    })

    // DELETE
    Register(http.MethodDelete, ApiPath + "/{id}", func (db *database.DB) http.HandlerFunc {
        return func (w http.ResponseWriter, r *http.Request){
            fmt.Fprintf(w, "Hello World " + r.PathValue("id") + "\n")
        }
    })
}
```

### The Conclusion
As you can see, we can register multiple routes with different http methods and handlers, they can share data between them such as the `ApiPath` we defined earlier.    
This allows us to both categorize and separate our routes neatly and logically into structures that make sense.    

This Plugin Style of registering routes, allows us to pretty much drag and drop files to activate/disable specific features of our server.   


