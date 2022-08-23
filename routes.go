package main

import (
	"github.com/audryhsu/plantopedia/models"
	"net/http"
	"strconv"
	"encoding/json"
	"log"
	"fmt"
	"github.com/gorilla/mux"
)
// centralized place for all routing is easier for maintainability, glanceability, and debug-ability
func (s *Server) handleRequests() {
	// Declare static file directory and point it assets dir
	// Handler routes reqs to respective filename, stripped "/assets/"" prefix to get correct filename 
	staticFileDirectory := http.Dir("./assets")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
		// match all routes starting with "/assets/" instead of absolute route
	s.Router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	s.Router.HandleFunc("/plants", s.handlePlants).Methods("GET")
	s.Router.HandleFunc("/plants/{id}", s.handlePlant).Methods("GET")
	s.Router.HandleFunc("/plants", s.handleCreatePlant).Methods("POST")
	s.Router.HandleFunc("/plants/{id}", s.handleDeletePlant).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", s.Router))

}

// handlers hang off the server and can access the dependencies via the s Server variable

func (s *Server) handlePlants(w http.ResponseWriter, r *http.Request) {
	var plants []models.Plant
	if result := s.Db.Find(&plants); result.Error != nil {
		log.Printf("Error getting all plants: %v\n", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("plants returned: {}", plants)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plants)
}

func (s *Server) handleCreatePlant(w http.ResponseWriter, r *http.Request) {
	newPlant := models.Plant{}
	// parse HTML form data
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form")
	}
	// get info about plant
	newPlant.Species = r.Form.Get("species")
	newPlant.Description = r.Form.Get("description")
	waterStr := r.Form.Get("water")
	newPlant.Water, _ = strconv.Atoi(waterStr)
	
	if result := s.Db.Create(&newPlant); result.Error != nil {
		log.Printf("Error creating plant %v\n", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("new plant created: %v\n", newPlant)
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPlant)

	// redirect user to index.html
	http.Redirect(w, r, "/assets/", http.StatusFound)
}

func (s *Server) handleDeletePlant(w http.ResponseWriter, r *http.Request) {
	// read the dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	w.Header().Add("Content-Type", "applicaton/json")
	encoder := json.NewEncoder(w) 
	var payload any
	var plant models.Plant

	// DELETE FROM plants WHERE id = id; does NOT return deleted row
	if result := s.Db.Delete(&plant, id); result.Error != nil {
		log.Println("error deleting plant")
		w.WriteHeader(http.StatusInternalServerError)

		} else if result.RowsAffected == 0 {
		log.Printf("Plant with id %v does not exist\n", id)
		w.WriteHeader(http.StatusNotFound)
		payload = fmt.Sprintf("No plant with id %v found.", id)
		
	} else {
		log.Printf("Deleted plant id: %v\n", id)
		w.WriteHeader(http.StatusOK)
	}

	encoder.Encode(payload)
}

func (s *Server) handlePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "applicaton/json")

	// find plant by id
	var plant models.Plant
	if result := s.Db.First(&plant, id); result.Error != nil {
		log.Printf("Error finding plant with id: %v: %v\n", id, result.Error)
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode("Plant not found.")
		return
	}

	log.Println(plant)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(plant)
}

// optional (not shown): write handlers to return a http.HandlerFunc instead to create a closure. helpful for one-time per handler initialization or if handler needs arguments 