package storage

import "github.com/satyampawar/go-project/internal/types"

type Storage interface {
	CreateStudent (name string , email string , age int ) (int64 , error)

	GetStudentByID(id int64) (types.Student, error)

	GetAllStudents() ([]types.Student, error)

}
