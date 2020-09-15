package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (server *Server) Initialize(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName string) {
	var err error

	if dbDriver == "mysql" {
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
		server.DB, err = gorm.Open(dbDriver, dbUrl)

		if err != nil {
			fmt.Printf("Cannot connect to %s database", dbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", dbDriver)
		}
	}

	if dbDriver == "postgres" {
		dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
		server.DB, err = gorm.Open(dbDriver, dbUrl)

		if err != nil {
			fmt.Printf("Cannot connect to %s database", dbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", dbDriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(address string) {
	fmt.Println("Listening to port 8001")
	log.Fatal(http.ListenAndServe(address, server.Router))
}
