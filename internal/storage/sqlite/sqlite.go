package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite" // pure Go driver (no CGO needed)

	"github.com/soumik171/Students-API/internal/config"
)

type Sqlite struct {
	Db *sql.DB //db connect
}

// create instance of struct: Have to use the func name as New
func New(cfg *config.Config) (*Sqlite, error) { //return instance of Sqlite and error

	db, err := sql.Open("sqlite", cfg.Storage_Path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err // return nil as nothing is inside at first
	}

	return &Sqlite{
		Db: db,
	}, nil

}
