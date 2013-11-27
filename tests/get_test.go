// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	// "allochi/inflect"
	"database/sql"
	"encoding/json"
	// "fmt"
	_ "github.com/lib/pq"
	// "strings"
	"testing"
	// "time"
)

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

}

func TestGet(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	var contacts []Contact
	ids := []int64{4, 5, 6, 7, 8, 9}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	err := Contactizer.Get(&contacts, request)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	var expected int
	db.QueryRow("select count(*) from json_contacts where id in (4, 5, 6, 7, 8, 9);").Scan(&expected)

	if len(contacts) != expected {
		t.Errorf("Expected %d contacts got %d", expected, len(contacts))
	}

}

func TestGetById(t *testing.T) {

	// Get few objects
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	var contacts []Contact
	ids := []int64{4, 5, 6, 7, 8, 9}
	err := Contactizer.Get(&contacts, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	var expected int
	db.QueryRow("select count(*) from json_contacts where id in (4, 5, 6, 7, 8, 9);").Scan(&expected)

	if len(contacts) != expected {
		t.Errorf("Expected %d contacts got %d", expected, len(contacts))
	}

}

func TestGetAll(t *testing.T) {

	// Get All objects
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	var contacts []Contact
	err := Contactizer.GetAll(&contacts)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	var expected int
	db.QueryRow("select count(*) from json_contacts;").Scan(&expected)

	if len(contacts) != expected {
		t.Errorf("Expected %d contacts got %d", expected, len(contacts))
	}

}

func TestGetRaw(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	// var contacts string
	ids := []int64{4, 5, 6, 7, 8, 9}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	data, err := Contactizer.GetRaw(Contact{}, request)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	expected := len(ids)

	var results interface{}
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		t.Errorf("Couldn't unpack json")
	}

	contacts := results.([]interface{})
	// got := strings.Count(data, "first_name")
	got := len(contacts)

	if expected != got {
		t.Errorf("Expected %d contacts got %d", expected, got)
	}

}
