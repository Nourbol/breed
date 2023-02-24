package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Permissions PermissionModel
	Breeds      interface {
		Insert(breed *Breed) error
		Get(id int64) (*Breed, error)
		Update(breed *Breed) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Breed, Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
		DeleteAllForUser(scope string, userID int64) error
	}
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Breeds:      BreedModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
