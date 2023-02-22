package data

import (
	"github.com/Nourbol/breed/internal/validator"
	"time"
)

type Breed struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	AvgCost     Cost      `json:"avg_cost,omitempty"`
	Countries   []string  `json:"countries,omitempty"`
	Version     int32     `json:"version"`
}

func ValidateBreed(v *validator.Validator, breed *Breed) {
	v.Check(breed.Name != "", "name", "must be provided")
	v.Check(len(breed.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(breed.Description != "", "descriptions", "must be provided")
	v.Check(len(breed.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(breed.AvgCost != 0, "avg_cost", "must be provided")
	v.Check(breed.AvgCost > 0, "avg_cost", "must be a positive integer")
	v.Check(len(breed.Countries) >= 1, "countries", "must contain at least 1 country")
	v.Check(len(breed.Countries) <= 5, "countries", "must not contain more than 5 countries")
	v.Check(validator.Unique(breed.Countries), "countries", "must not contain duplicate values")
}
