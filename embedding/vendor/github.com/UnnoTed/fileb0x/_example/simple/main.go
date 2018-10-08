//go:generate fileb0x b0x.yaml
package main

import (
	"log"
	"net/http"

	"github.com/UnnoTed/fileb0x/_example/simple/static"
)

func main() {
	files, err := static.WalkDirs("", false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ALL FILES", files)

	b, err := static.ReadFile("public/README.md")
	if err != nil {
		log.Fatal(err)
	}

	_ = b
	//log.Println(string(b))
	log.Println("try it -> http://localhost:8080/public/secrets.txt")

	// false = file system
	// true = handler
	as := false

	// try it -> http://localhost:8080/public/secrets.txt
	if as {
		// as Handler
		panic(http.ListenAndServe(":8080", static.Handler))
	} else {
		// as File System
		panic(http.ListenAndServe(":8080", http.FileServer(static.HTTP)))
	}
}
