package steps

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

func AddTaskWM(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Task{Error: err.Error()})
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Task{Error: err.Error()})
		return
	}

	if task.Title == "" || task.Title == " " {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Task{Error: "не указан заголовок задачи"})
		return
	}

	if task.Date == "" || task.Date == " " {
		task.Date = time.Now().Format(DateForFormat)
	}

	if task.Date == "today" {
		task.Date = time.Now().Format(DateForFormat)
	}

	parseDate, err := time.Parse(DateForFormat, task.Date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Task{Error: err.Error()})
		return
	}

	if parseDate.Before(time.Now()) {
		if task.Repeat == "" || task.Repeat == " " || task.Date == time.Now().Format(DateForFormat) {
			task.Date = time.Now().Format(DateForFormat)
		} else {
			task.Date, err = NextDateTask(time.Now(), task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				writeInfo(w, Task{Error: "ошибка функции вычисления даты выполнения задачи"})
				return
			}
		}
	} else {
		task.Date = task.Date
	}

	insertId, err := Insert(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Task{Error: "ошибка функции добавления записи в БД"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeInfo(w, Task{ID: insertId})

}

func writeInfo(w http.ResponseWriter, out any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(out)
}

func Insert(task Task) (int64, error) {
	row, err := DBConn.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}
