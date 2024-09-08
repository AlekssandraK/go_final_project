package main

import (
	"go_final_project/steps"
	"log"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env fle")
	}
	dbConn, err := steps.InitDB()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()
	steps.StartServer(dbConn)
}
