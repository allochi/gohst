package models

import (
	"strings"
	"time"
)

type Contact struct {
	Id             int64     `json:"-"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Title          string    `json:"title"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	PostalCode     string    `json:"postal_code"`
	Address        string    `json:"address_text"`
	Telephone      string    `json:"telephone"`
	Mobile         string    `json:"mobile"`
	Fax            string    `json:"fax"`
	Emails         []string  `json:"emails"`
	JobTitle       string    `json:"job_title"`
	Organization   string    `json:"organization"`
	IsOrganization bool      `json:"is_organization"`
	IsFocalPoint   bool      `json:"is_focal_point"`
	Department     string    `json:"department"`
	Role           string    `json:"role"`
	Comments       string    `json:"comments"`
	Categories     []string  `json:"categories"`
	Sectors        []string  `json:"sectors"`
	Interests      []string  `json:"interests"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (c *Contact) Name() string {
	if c.IsOrganization {
		return c.Organization
	}
	chunks := []string{c.Title, c.FirstName, c.LastName}
	return strings.TrimSpace(strings.Join(chunks, " "))
}
