package models

import (
	"strings"
	"time"
)

type Contact struct {
	Id           int64     `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Title        string    `json:"title"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	PostalCode   string    `json:"postal_code"`
	Telephone    string    `json:"telephone"`
	Mobile       string    `json:"mobile"`
	Fax          string    `json:"fax"`
	JobTitle     string    `json:"job_title"`
	Organization string    `json:"organization"`
	Department   string    `json:"department"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (c *Contact) Name() string {
	chunks := []string{c.Title, c.FirstName, c.LastName}
	return strings.TrimSpace(strings.Join(chunks, " "))
}
