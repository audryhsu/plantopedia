package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	
	// Declare static file directory and point it assets dir
	staticFileDirectory := http.Dir("./assets")
	// Handler routes reqs to respective filename, stripped "/assets/"" prefix to get correct filename 
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	
	// match all routes starting with "/assets/" instead of absolute route
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/plants", getPlantsHandler).Methods("GET")
	r.HandleFunc("/plants", createPlantHanlder).Methods("POST")

	return r
}

func main() {
	r:= newRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "hello world")
}