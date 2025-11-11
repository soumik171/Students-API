package storage

import "github.com/soumik171/Students-API/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudentsList() ([]types.Student, error)
	UpdateStudentInfo(id int64, student types.Student) (types.Student, error)
	DeleteStudent(id int64) error
}
