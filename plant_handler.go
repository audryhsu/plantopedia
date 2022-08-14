package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Plant struct {
	Species string `json:"species"`
	Description string `json:"description"`
}

var plants []Plant

func getPlantsHandler(w http.ResponseWriter, r *http.Request) {
	// convert plants var to json
	plantsBytes, err := json.Marshal(plants)
	// check for error and return resp to user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write JSON list of plants to the response object
	w.Write(plantsBytes)
}

func createPlantHanlder(w http.ResponseWriter, r *http.Request) {
	// create a new instance of a plant
	plant := Plant{}

	// parse HTML form data
	err := r.ParseForm()
	// error checking
	if err != nil {
		fmt.Println("Error parsing form")
	}
	// get info about plant and append to slice of plants 
	plant.Species = r.Form.Get("species")
	plant.Description = r.Form.Get("description")

	plants = append(plants, plant)

	// redirect user to index.html
	http.Redirect(w,r, "/assets/", http.StatusFound)
}