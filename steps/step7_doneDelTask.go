package steps

import (
	"database/sql"
	"net/http"
)

func TaskDone(w http.ResponseWriter, r *http.Request) {

	r.Method = http.MethodPost
	id := r.FormValue("id")
	task, err := ScanId(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	if task.Repeat == "" || task.Repeat == " " {
		err = DeleteId(id)

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

func DeleteId(id string) error {
	_, err := DBConn.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	err := DBConn.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}
	defer DBConn.Close()

	id := r.FormValue("id")
	_, err = ScanId(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	err = DeleteId(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "ошибка функции удаления записи БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeInfo(w, Task{})

}
