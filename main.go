package main

import "github.com/Abhishekkarunakaran/gyper/gyper"

func main() {
	g := gyper.New()
	g.Start("localhost","8888")
}
