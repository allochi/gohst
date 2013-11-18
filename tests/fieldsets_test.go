// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	"testing"
	"time"
)

type ContactType struct {
	Id        int64  `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string
	Addresses []Address   `json:"Addresses"`
	Emails    []string    `json:"emails"`
	Logs      []time.Time `json:"logs"`
	CreatedAt time.Time   `json:"-"`
}

type Address struct {
	Street  string
	Zipcode string
}

func TestScanFields(t *testing.T) {

	contact := ContactType{}

	tests := []struct {
		expected string
		got      string
	}{
		{"string", gohst.TypeOf(contact, "FirstName")},
		{"[]main.Address", gohst.TypeOf(contact, "Addresses")},
		{"[]time.Time", gohst.TypeOf(contact, "Logs")},
		{"time.Time", gohst.TypeOf(contact, "CreatedAt")},
	}

	for _, value := range tests {
		if value.got != value.expected {
			t.Errorf("ScanFields: Expected %d contacts got %d", value.expected, value.got)
		}
	}

}

func TestTypeName(t *testing.T) {

	tests := []struct {
		object   interface{}
		expected string
	}{
		{ContactType{}, "ContactType"},
		{&[]ContactType{}, "ContactType"},
		{[]ContactType{}, "ContactType"},
	}

	for _, test := range tests {
		got, _, _ := gohst.TypeName(test.object)
		if got != test.expected {
			t.Errorf("TypeName: Expected %d got %d", test.expected, got)
		}
	}

}
