package main

import (
	"ipManager/config"
	"ipManager/routers"
	"log"
	"net/http"
)

func main() {
	routers.Init()
	config.Init()

	err := http.ListenAndServe(":8080", routers.Mux)
	if err != nil {
		log.Fatal(err)
	}
}
