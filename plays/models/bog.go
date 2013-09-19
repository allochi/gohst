package models

import (
	"time"
)

type Bog struct {
	Id        int64 `json:"-"`
	Name      string
	Messages  []string
	Tags      []string
	Link      string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
