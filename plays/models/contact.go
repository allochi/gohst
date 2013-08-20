package models

import (
	"strings"
	"time"
)

type Contact struct {
	Id        int64     `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (c *Contact) Name() string {
	chunks := []string{c.Title, c.FirstName, c.LastName}
	return strings.TrimSpace(strings.Join(chunks, " "))
}
