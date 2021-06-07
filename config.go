package goa

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// Config ...
type Config struct {
	Service struct {
		Host   string `json:"host"`
		Port   string `json:"port"`
		Path   string `json:"path"`
		Method string `json:"method"`
	} `json:"service"`
	Cors struct {
		AllowOrigins string `json:"allow_origins"`
		AllowHeaders string `json:"allow_headers"`
	} `json:"cors"`
	Sort struct {
		DefaultSortOrder string   `json:"default_sort_order"`
		DefaultSortBy    string   `json:"default_sort_by"`
		ValidSortBy      []string `json:"valid_sort_by"`
	} `json:"sorting"`
	Pagination struct {
		DefaultPerPage string `json:"default_per_page"`
		MaxPerPage     string `json:"max_per_page"`
	} `json:"pagination"`
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"database"`
}

// ReadConfig ...
func ReadConfig() (config *Config) {
	env := os.Getenv("APP_ENV")
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if env == "" {
		log.Fatal("Missing Envionment!")
	}
	filePath := "./config/" + env + ".json"

	configFile, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer configFile.Close()

	json.NewDecoder(configFile).Decode(&config)

	config.Service.Path = strings.ToLower(config.Service.Path)
	config.Service.Method = strings.ToLower(config.Service.Method)

	return
}
