package data

import "time"

type Breed struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	AvgCost     Cost      `json:"avg_cost,omitempty"`
	Version     int32     `json:"version"`
}
