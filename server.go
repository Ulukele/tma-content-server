package main

import (
	"github.com/go-playground/validator/v10"
	"log"
)

// Validator
var validate = validator.New()

type Server struct {
	contentDBEngine *DBEngine
}

func NewServer() (*Server, error) {
	// configure content db engine
	// from environment
	//cDBC := DBConfig{
	//	Host:     os.Getenv("POSTGRES_HOST"),
	//	User:     os.Getenv("POSTGRES_USER"),
	//	Password: os.Getenv("POSTGRES_PASSWORD"),
	//	Name:     os.Getenv("POSTGRES_NAME"),
	//	Port:     os.Getenv("POSTGRES_PORT"),
	//	SSLMode:  "disable",
	//	Tz:       os.Getenv("POSTGRES_TZ"),
	//}
	cDBC := DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Name:     "postgres",
		Port:     "5432",
		SSLMode:  "disable",
		Tz:       "Asia/Novosibirsk",
	}

	log.Print("Try to create content db engine")
	contentEngine, err := NewDBEngine(cDBC)
	if err != nil {
		return nil, err
	}
	log.Print("Create content db engine")

	s := &Server{}
	s.contentDBEngine = contentEngine

	// init tables with models
	err = s.contentDBEngine.initTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}
