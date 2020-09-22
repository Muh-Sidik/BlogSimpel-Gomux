package api

import (
	"fmt"
	"log"
	"os"

	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/controllers"
	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/seeder"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	var err error

	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	server.Initialize(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName)

	seeder.Load(server.DB)

	server.Run(":8002")
}
