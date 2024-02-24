package main

import (
	"log"

	"github.com/royhq/go-play-app/shared/infra/bootstrap"
)

func main() {
	app, err := bootstrap.NewMainApp()
	if err != nil {
		log.Fatal("app bootstrap error:", err)
	}

	defer app.Shutdown()

	log.Fatal(app.ListenAndServe(":8080"))
}
