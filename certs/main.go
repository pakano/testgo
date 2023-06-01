package main

import "net/http"

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("123"))
	})
	http.ListenAndServeTLS(":8080", "./server.crt", "./server.key", nil)
}
