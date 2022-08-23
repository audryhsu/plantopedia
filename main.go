package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/audryhsu/plantopedia/db"
	"gorm.io/gorm"
)
// create a server struct that represents the service and holds all of its dependencies. no global state/variables! 
type Server struct {
	Router *mux.Router
	Db *gorm.DB
}

func main() {
	database := db.SetupDatabase()
	s := Server{
		Db: database,
		Router: mux.NewRouter().StrictSlash(true),
	}

	fmt.Println("Server started on 8080!")

	s.handleRequests()
}
