package steps

import (
	"database/sql"
	"net/http"
)

func TaskDone(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}

	err = db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}
	defer db.Close()

	r.Method = http.MethodPost
	id := r.FormValue("id")
	task, err := ScanId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	if task.Repeat == "" || task.Repeat == " " {
		err = DeleteId(db, id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeInfo(w, Err{Error: "ошибка функции удаления записи БД"})
			return
		}

		w.WriteHeader(http.StatusOK)
		writeInfo(w, Task{})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeInfo(w, Task{})
}

func DeleteId(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}

	err = db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}
	defer db.Close()

	id := r.FormValue("id")
	_, err = ScanId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	err = DeleteId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "ошибка функции удаления записи БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeInfo(w, Task{})

}
