# GoLang CheatSheet

---

## Commandline

### Core Commands
`go run [file.go]` compiles and runs a go file immediately    
`go build` compiles the package in the current directory and its dependencies into an executable without installing it    
`go install` compiles and installs packages into `$GOPATH/bin`    
`go fmt` automatically formats your source code according to the standard Go style    
`go test` runs test in the current package or specified packages    
`go get [module]` downloads and installs the specified module and its dependencies    
`go version` prints the installed go version    
`go env` prints information about the go environment variables    

### Module Management
`go mod init [path]` initializes a new module    
`go mod tidy` removes unused dependencies and adds missing ones to `go.mod`    
`go mod vendor` creates a vendor directory containing copies of all dependencies    
`go mod verify` checks that local dependencies match the expected hashes    
`go list -m all` lists the current module and all its dependencies    

### Documentation
`go doc [package/function]` displays information about a package or function    
`go help [command]` provides help for any Go subcommand    
