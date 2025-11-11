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

// Implement the interface of Creating student:

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

// Implement the interface of getting Student by Id :

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

// Implement the interface of Student List:

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

// Implement the interface of updating student:

func (s *Sqlite) UpdateStudentInfo(id int64, student types.Student) (types.Student, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name=?, email=?, age=? WHERE id=?")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(student.Name, student.Email, student.Age, id)
	if err != nil {
		return types.Student{}, err
	}

	RowsAffected, err := result.RowsAffected()

	if err != nil {
		return types.Student{}, err
	}
	if RowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", id)
	}

	var UpdateStudent types.Student

	// Using select query to diaplay the updated data, if don't want to update data, then no need to use the select query, just print the "process successed" like that

	if err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id=?", id).Scan(&UpdateStudent.Id, &UpdateStudent.Name, &UpdateStudent.Email, &UpdateStudent.Age); err != nil {
		return types.Student{}, err
	}

	return UpdateStudent, nil

}

// Implement the interface of deleting student:

func (s *Sqlite) DeleteStudent(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no student found with id %d", id)
	}

	return nil
}
