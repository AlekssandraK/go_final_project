package steps

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = "7540"

var dbConn *sql.DB

func StartServer(db *sql.DB) {
	dbConn = db
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("/api/nextdate", NextDate)
	mux.HandleFunc("/api/signin", auth)
	mux.HandleFunc("/api/task", authTask(selectFunc))
	mux.HandleFunc("/api/tasks/", authTask(searchHandler))
	mux.HandleFunc("/api/task/done", authTask(TaskDone))

	portStr, exists := os.LookupEnv("TODO_PORT")
	var currPort string
	if exists {
		currPort = portStr
	} else {
		currPort = port
	}
	fmt.Printf("Прослушивание порта: %s", currPort)
	err := http.ListenAndServe((":" + currPort), mux)
	if err != nil {
		log.Fatal(err)
	}

}

func selectFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskWM(w, r, dbConn)
		return
	case http.MethodGet:
		GetTaskId(w, r, dbConn)
		return
	case http.MethodPut:
		EditTask(w, r, dbConn)
		return
	case http.MethodDelete:
		DeleteTask(w, r, dbConn)
		return
	}
}
