package main

import (
	"log"

	"{{module}}/internal/api"
)

var version = "dev"

func main() {
	if err := api.Run(version); err != nil {
		log.Fatal(err)
	}
}
