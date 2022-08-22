package main

import (
	"fmt"
	"log"
	"net/http"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"encoding/json"

	// "database/sql"
	// "log"
	// _"github.com/jackc/pgx/v4/stdlib"
	"github.com/gorilla/mux"
)

var db *gorm.DB
var err error
var plants []Plant

type Plant struct {
	gorm.Model
	Id string `json:"Id"`
	Species     string `json:"species"`
	Description string `json:"description"`
	// Water int `json: water` 
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
	r.HandleFunc("/plants", createPlantHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func main() {
	fmt.Println("server started on 8080!")
	InitialMigration()

	// dummy plant data
	plants = []Plant{
		{Id: "1", Species: "Epipremnum aureum", Description: "Also known as golden pothos, it is a hardy vine plant from French Polynesia."},
		{Id: "2", Species: "Ficus lyrata", Description: "Also known as the fiddle-leaf fig, it is in the mulberry and fig family Moraceae. Native to western Africa."},
	}

	handleRequests()
}

// create table sql statements for us
func InitialMigration() {
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
		fmt.Println("Created plants table successfully!")
	}
}

func getPlantsHandler(w http.ResponseWriter, r *http.Request) {
	// url is username:password@localhost:5432/DATABASE_NAME
	dbURL:= "postgres://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	// check for error and return resp to user
	if err != nil {
		fmt.Println(fmt.Errorf("Error connecting to postgres: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []Plant
	db.Find(&plants)
	fmt.Println("plants returned: {}", result)

	json.NewEncoder(w).Encode(result)

	fmt.Printf("sent data: %v\n", result)

	// write JSON list of plants to the response object
	// w.Write(plantsBytes)
}

func createPlantHandler(w http.ResponseWriter, r *http.Request) {
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

	// update global variable with new plant 
	plants = append(plants, plant)
	json.NewEncoder(w).Encode(plants)
	fmt.Printf("new plant added: %v\n", plants)

	// redirect user to index.html
	http.Redirect(w, r, "/assets/", http.StatusFound)
}