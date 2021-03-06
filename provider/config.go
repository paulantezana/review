package provider

import (
	"encoding/json"
	"log"
	"os"
)

// Config models
type Database struct {
	Server   string
	Port     string
	User     string
	Pass     string
	Database string
}

type Email struct {
	From     string
	Password string
	Server   string
	Host     string
}

type Global struct {
	Paginate uint
}

type Config struct {
	Database Database
	Email    Email
	Server   Server
	Global   Global
}

type Server struct {
	Port   string
	Host   string
	Socket string
	Key    string
}

// GetConfig return configuration from database json
func GetConfig() Config {
	var c Config

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
