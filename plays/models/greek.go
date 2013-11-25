package models

import (
	"time"
)

type Greek struct {
	Id        int64 `json:"-"`
	Name      string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
