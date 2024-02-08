package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}

		log.Println("POST received with body: ", string(body))
		answer := "Received POST request: " + string(body) + "\n"
		w.Write([]byte(answer))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/test", handlePostRequest)

	http.ListenAndServe(":8080", nil)
}
