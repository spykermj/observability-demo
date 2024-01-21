package main

import (
	"log"

	"spykerman.co.uk/roller/internal/serve"
)

func main() {
	log.Fatalln(serve.ServeDice())
}
