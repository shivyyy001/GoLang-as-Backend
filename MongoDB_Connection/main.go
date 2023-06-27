package main

import (
	"fmt"
	"github.com/shivansh/mongodbconnection/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Listerning to port...")
	log.Fatal(http.ListenAndServe(":4009", r))
}
