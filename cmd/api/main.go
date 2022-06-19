package main

import (
	"log"

	"github.com/Palantay/constanta/internal/app/api"
)

var ()

func init() {

}

func main() {
	config := api.NewConfig()

	server := api.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
