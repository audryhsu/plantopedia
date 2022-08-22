package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "database/sql"
	// "log"
	// _"github.com/jackc/pgx/v4/stdlib"
	"github.com/gorilla/mux"
)

var db *gorm.DB
var err error

type Plant struct {
	gorm.Model
	Id int `json:"id" gorm:"primaryKey"`
	Species     string `json:"species"`
	Description string `json:"description"`
	Water int `json:"water"`
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true) 
	
	// Declare static file directory and point it assets dir
	staticFileDirectory := http.Dir("./assets")
	// Handler routes reqs to respective filename, stripped "/assets/"" prefix to get correct filename 
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
		// match all routes starting with "/assets/" instead of absolute route
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/plants", getPlantsHandler).Methods("GET")
	r.HandleFunc("/plants/{id}", getPlantHandler).Methods("GET")
	r.HandleFunc("/plants", createPlantHandler).Methods("POST")
	r.HandleFunc("/plants/{id}", deletePlantHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func main() {
	fmt.Println("server started on 8080!")
	InitialMigration()
	seedDatabase()

	handleRequests()
}

// create table sql statements for us
func InitialMigration() {
	
	// set global db var	
	// url is username:password@localhost:5432/DATABASE_NAME
	// dbURL:= "postgres://postgres:postgres@localhost:5432/postgres"
	// db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
	}

	err = db.AutoMigrate(&Plant{})
	if err != nil {
		log.Fatalf("Couldnt migrate table\n")
	} else {
		log.Println("Created plants table successfully!")
	}

}

func getPlantsHandler(w http.ResponseWriter, r *http.Request) {
	// url is username:password@localhost:5432/DATABASE_NAME
	// dbURL:= "postgres://postgres:postgres@localhost:5432/postgres"
	// db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	
	var plants []Plant
	if result := db.Find(&plants); result.Error != nil {
		log.Printf("Error getting all plants: %v\n", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("plants returned: {}", plants)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plants)
}

func createPlantHandler(w http.ResponseWriter, r *http.Request) {
	// create a new instance of a plant
	newPlant := Plant{}
	// parse HTML form data
	err := r.ParseForm()
	// error checking
	if err != nil {
		log.Println("Error parsing form")
	}
	// get info about plant
	newPlant.Species = r.Form.Get("species")
	newPlant.Description = r.Form.Get("description")
	waterStr := r.Form.Get("water")
	newPlant.Water, _ = strconv.Atoi(waterStr)
	
	// newPlant := Plant{ Species: "Zamioculcas zamiifolia", Description: "Also known as the ZZ plant, it is a tropical perennial plant native to eastern Africa, from southern Kenya to northeastern South Africa. It's low-light tolerant and has rhizomes that can retain water. Water less frequently.", Water: 1}
	
	if result := db.Create(&newPlant); result.Error != nil {
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

func deletePlantHandler(w http.ResponseWriter, r *http.Request) {
	// read the dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	w.Header().Add("Content-Type", "applicaton/json")
	encoder := json.NewEncoder(w) 
	var payload any
	var plant Plant

	// DELETE FROM plants WHERE id = id
	if result := db.Delete(&plant, id); result.Error != nil {
		log.Println("error deleting plant")
		w.WriteHeader(http.StatusInternalServerError)

		} else if result.RowsAffected == 0 {
		log.Printf("Plant with id %v does not exist\n", id)
		w.WriteHeader(http.StatusNotFound)
		payload = fmt.Sprintf("No plant with id %v found.", id)
		
	} else {
		log.Printf("Deleted plant id: %v\n", id)
		w.Header().Add("Content-Type", "applicaton/json")
		w.WriteHeader(http.StatusOK)
		payload = plant
	}

	encoder.Encode(payload)
}

func getPlantHandler(w http.ResponseWriter, r *http.Request) {
	// read the dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "applicaton/json")

	// find plant by id
	var plant Plant
	if result := db.First(&plant, id); result.Error != nil {
		log.Printf("Error finding plant with id: %v: %v\n", id, result.Error)
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode("Plant not found.")
		return
	}

	log.Println(plant)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(plant)
}

func seedDatabase() {
	// seed database if no rows in it.
	var existingPlants []Plant

	result := db.Find(&existingPlants)
	if result.RowsAffected == 0 {
		plants := []Plant{
		{Species: "Epipremnum aureum", Description: "Also known as golden pothos, it is a hardy vine plant from French Polynesia. Tolerates low light and likes vines to hang.", Water: 3},
		{Species: "Ficus lyrata", Description: "Also known as the fiddle-leaf fig, it is in the mulberry and fig family Moraceae. Native to western Africa.", Water: 2},
		{ Species: "Zamioculcas zamiifolia", Description: "Also known as the ZZ plant, it is a tropical perennial plant native to eastern Africa, from southern Kenya to northeastern South Africa. It's low-light tolerant and has rhizomes that can retain water. Water less frequently.", Water: 1},
	}
	
	if result := db.Create(&plants); result.Error!= nil {
		log.Printf("Error seeding db: %v", result.Error)
	}
	
} else {
	fmt.Println("Database already data, no seeding.")
}
	
}