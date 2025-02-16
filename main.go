package main

import (
	"fmt"
	"log"

	"github.com/Abhishekkarunakaran/gyper/gyper"
)

type function func(value string) error

type data struct {
	name       string
	handleFunc function
}

func main() {

	// d := data{
	// 	name: "function 1",
	// 	handleFunc: func(value string) error {
	// 		return fmt.Errorf("sample error")
	// 	},
	// }

	// run(d.handleFunc, d.name)
	g := gyper.New()
	if err := g.Start("localhost", "8888"); err != nil {
		log.Fatal(err.Error())
		g.Stop()
	}
}

func run(handle function, name string) {
	if err := handle(name); err != nil {
		fmt.Println(err)
	}
}
