# fproto

[![GoDoc](https://godoc.org/github.com/RangelReale/fproto?status.svg)](https://godoc.org/github.com/RangelReale/fproto)

Package for parsing Google Protocol Buffers into Go structs [.proto files version 2 + 3] (https://developers.google.com/protocol-buffers/docs/reference/proto3-spec)

See elements.go to see the structs definitions.

Uses [https://github.com/emicklei/proto](https://github.com/emicklei/proto) for parsing.

### install

    go get -u -v github.com/RangelReale/fproto

### usage

	package main

	import (
	    "fmt"
        "os"

        "github.com/RangelReale/fproto"
	)

	func main() {
        file, err := os.Open("/file/name")
        if err != nil {
            return nil, err
        }
        defer file.Close()
    
        protofile, err := fproto.Parse(file)
        if err != nil {
            return nil, err
        }
        
        fmt.Printf("Package name: %s\n", protofile.PackageName)
	}
	
### related
 
 * [https://github.com/RangelReale/fdep](https://github.com/RangelReale/fdep) Package for building relationships between 
 proto files and extracting types, helping creating source code generators.

 * [https://github.com/RangelReale/fproto-wrap](https://github.com/RangelReale/fproto-wrap)
    Generates easier-to-use wrappers for the standard protobuf generated files.
	
### author

Rangel Reale (rangelspam@gmail.com)
