package main

import (
	"fmt"
	"log"
	"net/http"

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

type Profile struct {
	Name string `json:"name" xml:"name"`
	Age  int    `json:"age" xml:"age"`
}

func (p *Profile) String() string {
	return fmt.Sprintf("\nname : %s, age : %d", p.Name, p.Age)
}

func run(g gyper.Context) {
	fmt.Printf("method : %s\n", g.Request.Method)
	fmt.Printf("path : %s\n", g.Request.Path)
	fmt.Printf("headers : %v", g.Request.Header)
	var profile Profile
	if err := g.Bind(&profile); err != nil {
		fmt.Println(err.Error())
	}

	g.XML(http.StatusOK,profile)
	fmt.Println(profile.String())
}
