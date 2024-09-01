package main

import (
	"github.com/Rafael24595/go-api-render/src/infrastructure"
	"log"
	"net/http"
)

func main() {
	infrastructure.NewController()
	log.Fatalln(http.ListenAndServe(":8080", nil))
}