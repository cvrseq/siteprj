package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./frontend/pages"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/pages/registration.html")
	})

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
