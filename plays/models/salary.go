package models

import (
	"time"
)

type Salary struct {
	Id        int64 `json:"-"`
	Name      string
	Amount    float32
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
