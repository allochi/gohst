package models

import (
	"time"
)

type MailingList struct {
	Id          int64     `json:"-"`
	UserId      int64     `json:"user_id"`
	Shared      bool      `json:"shared"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Template    string    `json:"letter_template"`
	ContactIds  []int64   `json:"contact_ids"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
