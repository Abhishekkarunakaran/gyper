package main

import (
	"log"

	"github.com/Abhishekkarunakaran/gyper/gyper"
)

func main() {
	g := gyper.New()
	if err := g.Start("localhost","8888"); err != nil {
		log.Fatal(err.Error())
		g.Stop()
	}
}
