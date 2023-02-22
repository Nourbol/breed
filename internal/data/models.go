package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Breeds interface {
		Insert(movie *Breed) error
		Get(id int64) (*Breed, error)
		Update(movie *Breed) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Breed, Metadata, error)
	}
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Breeds: BreedModel{DB: db},
	}
}
