![Hello World](images/001_golang-hello-world.png)

```go
// declare package as main
package main;

// import the fmt stdlib package for i/o
import "fmt";

// define the main function
func main() {

        // print hello world to the console
        fmt.Println("Hello World");

}
```
Go uses javascript style comments, e.g `// single line comment` and `/* multiline comments */`    
`package main` Every go file must start with a package declaration   
`import "fmt"` imports the fmt package that contains functions for formatted input output operations (such as printf)    
`func main(){}` defines the main function executed when run    
`fmt.Println("tect here")` prints a line of text to stdout    

