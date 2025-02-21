package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Abhishekkarunakaran/gyper/gyper"
)

func main() {
	g := gyper.New()
	g.GET("/v1/private/profile", getValue)
	g.POST("/v1/private/profile", save)
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

var profile Profile

func save(c gyper.Context) {
	if err := c.Bind(&profile); err != nil {
		fmt.Println(err.Error())
	}
	_ = c.JSON(http.StatusOK, profile)
	// fmt.Println(profile.String())
}

func getValue(c gyper.Context) {
	_ = c.JSON(http.StatusOK, profile)
}
