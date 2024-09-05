package main

import (
	"database/sql"
	"fmt"
	"go_final_project/steps"

	_ "modernc.org/sqlite"
)

func main() {
	success := steps.CreateDB()
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
	}
	if success {
		steps.StartServer(db)
	} else {
		fmt.Println("Не удалось создать БД")
	}

}
