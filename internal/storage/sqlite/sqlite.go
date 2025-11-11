package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // pure Go driver (no CGO needed)

	"github.com/soumik171/Students-API/internal/config"
	"github.com/soumik171/Students-API/internal/types"
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

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	// Declare the query:

	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age) VALUES(?,?,?)") // ? used to avoid sql injection

	if err != nil {
		return 0, err
	}

	defer stmt.Close() // close the query statement(Insert Into......)

	// Execute the query:

	result, err := stmt.Exec(name, email, age) // value find

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	//pass the ref, if no need to parse then use Query, if have to pass something, then have to use based on which attributes i wanna use like:QueryRow(),QueryContext()....

	errV := stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if errV != nil {
		if errV == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", errV)
	}

	return student, nil

}

func (s *Sqlite) GetStudentsList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student //as return list, create a container that hold the result

	// For list have to use the for loop
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil

}
