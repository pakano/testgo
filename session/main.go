package main

import (
	"net/http"
	"test/session/user"

	"github.com/gorilla/mux"
)

func registerRouter(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", user.Login).Methods("POST")
	userRouter.HandleFunc("/secret", user.Secret)
	userRouter.HandleFunc("/logout", user.LogOut)
}

func main() {
	r := mux.NewRouter()
	registerRouter(r)

	http.ListenAndServe(":9898", r)
}
