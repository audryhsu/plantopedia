package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	// "database/sql"
// 	"log"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var db *gorm.DB
// var err error
// var plants []Plant

// type Plant struct {
// 	gorm.Model
// 	Species     string `json:"species"`
// 	Description string `json:"description"`
// 	// Water int `json: water` 
// }

// // create table sql statements for us
// func InitialMigration() {
// 	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

// 	db, err = gorm.Open(postgres.New(postgres.Config{
// 		DSN:                  dsn,
// 		PreferSimpleProtocol: true,
// 	}), &gorm.Config{})

// 	if err != nil {
// 		log.Fatalf("Could not connect to db: %v\n", err)
// 	}

// 	err = db.AutoMigrate(&Plant{})
// 	if err != nil {
// 		log.Fatalf("Couldnt migrate table\n")
// 	} else {
// 		fmt.Println("Created plants table successfully!")
// 	}
// }

// func getPlantsHandler(w http.ResponseWriter, r *http.Request) {
// 	// convert plants var to json
// 	// plantsBytes, err := json.Marshal(plants)

// 	db, err := gorm.Open(postgres.Open("postgres"))
// 	// check for error and return resp to user
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("Error connecting to postgres: %v", err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	var plants []Plant
// 	db.Find(&plants)
// 	fmt.Println("plants returned: {}", plants)

// 	json.NewEncoder(w).Encode(plants)

// 	// write JSON list of plants to the response object
// 	// w.Write(plantsBytes)
// }

// func createPlantHandler(w http.ResponseWriter, r *http.Request) {
// 	// create a new instance of a plant
// 	plant := Plant{}

// 	// parse HTML form data
// 	err := r.ParseForm()
// 	// error checking
// 	if err != nil {
// 		fmt.Println("Error parsing form")
// 	}
// 	// get info about plant and append to slice of plants
// 	plant.Species = r.Form.Get("species")
// 	plant.Description = r.Form.Get("description")

// 	plants = append(plants, plant)

// 	// redirect user to index.html
// 	http.Redirect(w, r, "/assets/", http.StatusFound)
// }
