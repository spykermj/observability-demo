package main

import (
	"context"
	"log"

	"spykerman.co.uk/roller/internal/otel"
	"spykerman.co.uk/roller/internal/serve"
)

func main() {
	tp, err := otel.InitTrace(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := tp.Shutdown(context.Background())

		if err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	log.Fatalln(serve.ServeDice())
}
