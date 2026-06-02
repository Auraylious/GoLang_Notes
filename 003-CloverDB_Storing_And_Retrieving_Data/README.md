# CloverDB - Storing and Retrieving Data

![CloverDB Logo](https://raw.githubusercontent.com/ostafen/clover/v2/.github/logo-white.png#gh-dark-mode-only)    

[CloverDB](https://github.com/ostafen/clover) is a lightweight, No SQL, document oriented database written in Go. It is intended to be embedded directly into go projects without the need for maintaining a network server like mongodb.    

### Similar to MongoDB
CloverDB is similar to MongoDB   
- has schemaless collections or documents
- it adds a unique UUID `_id` field to each document
- NoSQL queries, can filter, sort, and extract unrelated data

### Installation
to use cloverdb we must first initialize a go module: `go mod init CloverTesting`
then we can install it with `go get github.com/ostafen/clover/v2`

### Opening the Database
for locally stored data, clover uses Bolt under the hood by default, but can also be configured to use badger for an in-memory database as well.    
when opening the database, clover will create the directory path if it does not exist.   

to open a local datastore we use `clover.Open("/Path/to/DB")`   
this returns a `clover.DB` object, which is our open connection to the database, and we will need to run `db.Close()` on it when were completely finished with it.

### Storing Data to the Database
CloverDB collections are entirely schemaless. You can seamlessly push maps alongside structs to the exact same collection without prior definitions    

First we create a collection with `db.CreateCollection("collection_name")`   

Then we must convert raw JSON bytes or Go structs into a `clover.Document` object.   
Then we can insert the document into the collection like so `db.InsertOne("collection", doc)`    

### Retrieving Data from the Database

to retrieve data, we create a query object that contains what were looking for, and then run the query    
An Example: `db.Find(query.NewQuery("users").Where(query.Field("age").Gt(25)))`    

- `db.Find` runs the query
- `query.NewQuery("collection")` creates a new query object for a certain collection
- `.Where()` filters our query results by certain criteria
- `query.Field("field")` field object to look for inside of documents
- `.Gt(25)` filter documents based on criteria (specifically if the value is greater than 25)

you can also retrieve all the documents in a collection like so: `db.FindAll(query.NewQuery("collection"))`    

we can iterate through them like so:
```go
docs, _ := db.FindAll(query.NewQuery("myCollection"))
or _, doc := range docs {
      log.Println(doc)
}
```


# Putting it all Together
[this code]("Example.go")    
```go
package main

// we directly import the document and query sub packages as well as the main clover package
import (
    "log"
    "encoding/json"
    "github.com/ostafen/clover/v2"
    "github.com/ostafen/clover/v2/document"
    "github.com/ostafen/clover/v2/query"
)

func main(){

    // open the database in a subdirectory of the current directly
    db, _ := clover.Open("clover-database")

    // close it when were done
    defer db.Close()

    // create our collection if it doesnt exist
    db.CreateCollection("myCollection")

    // create a new document object, manually creating fields and values
    doc := document.NewDocument()
    doc.Set("hello", "clover!")

    // insert the document to the collection, InsertOne returns the new documents _id
    docId, _ := db.InsertOne("myCollection", doc)
    log.Println("Inserted: " + docId)

    // convert raw json to a document object, by unmarshalling it to a map
    rawJSON := `{"name": "Alice Smith", "age": 28, "roles": ["admin", "user"]}`
    var jsonMap map[string]interface{}
    json.Unmarshal([]byte(rawJSON), &jsonMap)
    ourNewDocument := document.NewDocumentOf(jsonMap)

    // insert that into the db as well
    docId, _ = db.InsertOne("myCollection", ourNewDocument)
    log.Println("Inserted: " + docId)

	// Initializing a Clover Document directly from a struct
    type User struct {
        Name  string   `json:"name"`
        Age   int      `json:"age"`
        Roles []string `json:"roles"`
    }
   	userStruct := User{
		Name:  "Bob Jones",
		Age:   34,
		Roles: []string{"moderator"},
	}
	docFromStruct := document.NewDocumentOf(userStruct)

    // insert that into the db as well
    docId, _ = db.InsertOne("myCollection", docFromStruct)
    log.Println("Inserted: " + docId)

    // Finally lets retrieve all documents in our collection and display them to see how they look
    docs, _ := db.FindAll(query.NewQuery("myCollection"))
    for _, doc := range docs {
          log.Println(doc)
	  // lets print their _id field
	  log.Println(doc.Get("_id"))
	  // fields are case sensitive and returns nil when undefined
	  log.Println(doc.Get("name"))
    }

}
```
