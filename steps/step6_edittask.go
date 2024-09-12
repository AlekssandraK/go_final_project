package steps

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type TaskID struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTaskId(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := ScanId(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	w.WriteHeader(http.StatusOK)

	writeInfo(w, TaskID{ID: task.ID, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat})
}

func ScanId(id string) (TaskID, error) {
	row := DBConn.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	var task TaskID
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return task, err
	}

	return task, nil
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	var buff bytes.Buffer
	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}

	var task TaskID

	if err = json.Unmarshal(buff.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: err.Error()})
		return
	}

	if _, err := ScanId(strconv.Itoa(int(task.ID))); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "задача не найдена"})
		return
	}

	if task.Title == "" || task.Title == " " {
		w.WriteHeader(http.StatusBadRequest)
		writeInfo(w, Err{Error: "не указан заголовок задачи"})
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
		writeInfo(w, Err{Error: err.Error()})
		return
	}

	if parseDate.Before(time.Now()) {
		if task.Repeat == "" || task.Repeat == " " || task.Date == time.Now().Format(DateForFormat) {
			task.Date = time.Now().Format(DateForFormat)
		} else {
			task.Date, err = NextDateTask(time.Now(), task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				writeInfo(w, Err{Error: "ошибка функции вычисления даты выполнения задачи"})
				return
			}
		}
	}

	_, err = DBConn.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeInfo(w, Err{Error: "ошибка функции изменения записи в БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeInfo(w, Task{})

}
