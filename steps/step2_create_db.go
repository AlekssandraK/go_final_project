package steps

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() (dbConn *sql.DB, err error) {

	path, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbfile := os.Getenv("TODO_DBFILE")
	DBConn, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		panic(err)
	}

	db := filepath.Join(filepath.Dir(path), dbfile)
	_, err = os.Stat(db)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DBConn.Exec(`CREATE TABLE IF NOT EXISTS scheduler
		(id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL DEFAULT '',
		  title VARCHAR(128) NOT NULL DEFAULT '', 
		  comment VARCHAR(256) NOT NULL DEFAULT '',
		  repeat VARCHAR(128) NOT NULL DEFAULT '')`, `CREATE INDEX date_index ON scheduler(date)`)

	if err != nil {
		log.Fatal(err)
	}
	return DBConn, nil
}
