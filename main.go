package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "8080"

func main() {

	http.HandleFunc("/", handleDisplaySurface)

	log.Printf("Running on http://localhost:%s", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleDisplaySurface(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Display Surface")
}
