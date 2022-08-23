package db

import (
	"fmt"
	"github.com/audryhsu/plantopedia/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var err error

func SetupDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to DB: %v\n", err)
	}

	// create table sql statements for us
	err = db.AutoMigrate(&models.Plant{})
	if err != nil {
		log.Fatalf("Couldnt migrate table\n")
	} else {
		fmt.Println("Created plants table successfully!")
	}

	// seed database if no rows in it.
	seed(db)

	return db
}

func seed(db *gorm.DB) {
	var existingPlants []models.Plant

	result := db.Find(&existingPlants)
	if result.RowsAffected == 0 {
		plants := []models.Plant{
			{Species: "Epipremnum aureum", Description: "Also known as golden pothos, it is a hardy vine plant from French Polynesia. Tolerates low light and likes vines to hang.", Water: 3},
			{Species: "Ficus lyrata", Description: "Also known as the fiddle-leaf fig, it is in the mulberry and fig family Moraceae. Native to western Africa.", Water: 2},
			{Species: "Zamioculcas zamiifolia", Description: "Also known as the ZZ plant, it is a tropical perennial plant native to eastern Africa, from southern Kenya to northeastern South Africa. It's low-light tolerant and has rhizomes that can retain water. Water less frequently.", Water: 1},
		}

		if result := db.Create(&plants); result.Error != nil {
			log.Printf("Error seeding DB: %v", result.Error)
		} else {
			fmt.Println("Database already data, no seeding.")
		}
	}
}

// url is username:password@localhost:5432/DATABASE_NAME
// dbURL:= "postgres://postgres:postgres@localhost:5432/postgres"
// db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
