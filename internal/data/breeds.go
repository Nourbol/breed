package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nourbol/breed/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
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

type BreedModel struct {
	DB *pgxpool.Pool
}

func (m BreedModel) Insert(breed *Breed) error {
	query := `INSERT INTO breeds (name, description, avg_cost, countries)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	args := []any{breed.Name, breed.Description, breed.AvgCost, breed.Countries}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&breed.ID, &breed.CreatedAt, &breed.Version)
}

func (m BreedModel) Get(id int64) (*Breed, error) {
	query := `SELECT id, created_at, name, description, avg_cost, countries, version
			  FROM breeds
			  WHERE id = $1`

	var breed Breed

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRow(ctx, query, id).
		Scan(&breed.ID,
			&breed.CreatedAt,
			&breed.Name,
			&breed.Description,
			&breed.AvgCost,
			&breed.Countries,
			&breed.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &breed, nil
}

func (m BreedModel) Update(breed *Breed) error {
	query := `UPDATE breeds
			  SET name = $1, description = $2, avg_cost = $3, countries = $4, version = version + 1
			  WHERE id = $5 AND version = $6
			  RETURNING version`

	args := []any{
		breed.Name,
		breed.Description,
		breed.AvgCost,
		breed.Countries,
		breed.ID,
		breed.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRow(ctx, query, args...).Scan(&breed.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m BreedModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM breeds
			  WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m BreedModel) GetAll(name string, countries []string, filters Filters) ([]*Breed, Metadata, error) {
	query := fmt.Sprintf(`
			  SELECT COUNT(*) OVER (), id, created_at, name, description, avg_cost, countries, version 
			  FROM breeds 
			  WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			  AND (countries @> $2 OR $2 = '{}')
			  ORDER BY %s %s, id ASC
			  LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, countries, filters.limit(), filters.offset()}

	rows, err := m.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var breeds []*Breed

	for rows.Next() {
		var breed Breed

		if err := rows.Scan(
			&totalRecords,
			&breed.ID,
			&breed.CreatedAt,
			&breed.Name,
			&breed.Description,
			&breed.AvgCost,
			&breed.Countries,
			&breed.Version); err != nil {
			return nil, Metadata{}, err
		}

		breeds = append(breeds, &breed)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return breeds, metadata, nil
}
