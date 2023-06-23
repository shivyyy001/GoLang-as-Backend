package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// using gorilla mux

	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler)
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)

	fmt.Printf("Starting server at port 9000\n")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// how to parse form data ->
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	fmt.Println(r.Form)
	fmt.Println(r.PostForm)

	// how to get form values ->
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}
