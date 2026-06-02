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
