package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var database MessageStorage

func handleRequest(directory string, message string) {
	http.HandleFunc(directory, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != directory {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, message)
	})
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintln(w, "Visit the URL below to redeem the one-time message:")
	message := r.FormValue("message")

	host := r.Host
	id := database.Store(message)
	url := fmt.Sprintf("http://%s/view?id=%s", host, id)
	fmt.Fprintln(w, url)
}

func randomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	var url bytes.Buffer
	for i := 0; i < 10; i++ {
		url.WriteByte(charset[rand.Intn(len(charset))])
	}
	return url.String()
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")
	message, err := database.Redeem(id)
	if err != nil {
		http.NotFound(w, r)
		return 
	}
	fmt.Fprintf(w, "This is the secret message: %s", message)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		fileServer := http.FileServer(http.Dir("./static"))
		fileServer.ServeHTTP(w, r)
	} else if r.Method == "POST" {
		formHandler(w, r)
	}
}

func main() {
	database = NewInMemoryMessageStorage()

	http.HandleFunc("/", indexHandler)
	//handleRequest("/form", "Input your secret message below.")
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/view", viewHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
