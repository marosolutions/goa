package goapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceInterface interface {
	Controller(config *Config, w http.ResponseWriter, req *http.Request)
}

var service ServiceInterface
var DB *pgxpool.Pool

// StartService ...
func StartService(srv ServiceInterface) {
	service = srv

	// Read configuration from the file path
	config := ReadConfig()

	// Create a new Database client
	DB = config.newDBClient()
	defer DB.Close()

	http.HandleFunc(config.Service.Path, config.Router) // Load all the routes

	fmt.Println("Starting service:", config)

	addr := config.Service.Host + ":" + config.Service.Port
	err := http.ListenAndServe(addr, nil) // Start the API server
	if err != nil {
		log.Fatal("Error! server failed to start.", err)
	}
}
