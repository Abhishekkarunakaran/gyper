package main

import (
	"fmt"
	"log"

	"github.com/Abhishekkarunakaran/gyper/gyper"
)

func main() {
	g := gyper.New()
	g.GET("", run)
	g.POST("/v1/private", run)
	if err := g.Start("localhost", "8888"); err != nil {
		log.Fatal(err.Error())
		g.Stop()
	}
}

func run(g gyper.Context) {
	fmt.Printf("method : %s\n", g.Request.Method)
	fmt.Printf("path : %s\n", g.Request.Path)
	fmt.Printf("headers : %v", g.Request.Header)
}
