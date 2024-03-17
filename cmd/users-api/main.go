package main

import (
	"log"

	"github.com/royhq/go-play-app/shared/bootstrap"
)

func main() {
	app, err := bootstrap.NewUsersAPI()
	if err != nil {
		log.Fatal("app bootstrap error:", err)
	}

	defer app.Shutdown()

	log.Fatal(app.ListenAndServe(":8080"))
}
