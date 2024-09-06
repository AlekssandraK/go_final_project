package main

import (
	"go_final_project/steps"

	_ "modernc.org/sqlite"
)

func main() {
	dbConn, err := steps.InitDB()
	if err != nil {
		panic(err)
	}
	steps.StartServer(dbConn)
	dbConn.Close()
}
