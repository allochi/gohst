// test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

var Contactizer gohst.DataStore

func init() {

	AllochiContactizer := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	Contactizer.Register("Contactizer", AllochiContactizer)
	Contactizer.Connect()
	// defer Contactizer.Disconnect()
	// gohst.Register("Contactizer", Contactizer)
}

func TestReadData(t *testing.T) {

	var contacts []Contact

	var err error
	var expected int
	var got int

	ids := []int64{4, 5, 6, 7, 8, 9}
	expected = len(ids)
	err = Contactizer.GET(&contacts, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	got = len(contacts)
	if got != expected {
		t.Errorf("Get Few: Expected %d contacts got %d", expected, got)
	}

	var allContacts []Contact
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	db.QueryRow("select count(*) from json_contacts;").Scan(&expected)
	err = Contactizer.GET(&allContacts, []int64{})
	got = len(allContacts)
	if got != expected {
		t.Errorf("Get All: Expected %d contacts got %d", expected, got)
	}
}

func TestInsertAndDelete(t *testing.T) {

	bog := Bog{}
	bog.Name = "Ali Anwar"
	bog.Messages = []string{"Go is cool", "PostgreSQL is cool"}
	bog.Tags = []string{"Go", "PostgreSQL"}
	bog.Link = "http://www.allochi.com"

	// Return and ID
	err := Contactizer.PUT(&bog)

	if err != nil {
		t.Fatalf("gohst error: %s", err)
	}

	if bog.Id <= 0 {
		t.Errorf("PUT: Expected an ID > 0 got %d", bog.Id)
	}

	ids := []int64{bog.Id}

	err = Contactizer.DELETE([]Bog{}, ids)

	if err != nil {
		t.Fatalf("gohst error: %s", err)
	}

	var bogs []Bog
	err = Contactizer.GET(&bogs, ids)

	if len(bogs) > 0 {
		t.Errorf("DELETE: Didn't expect an object got %d", len(bogs))
	}

}

// 17: BenchmarkReadData	    5000	    347401 ns/op
// 13: BenchmarkReadData	    5000	    348989 ns/op
// 13: BenchmarkReadData	    5000	    361726 ns/op
// 13: BenchmarkReadData	   10000	    342638 ns/op (2918 op/s)
// 13: BenchmarkReadData	   10000	    313696 ns/op (3187 op/s)
func BenchmarkReadData(b *testing.B) {

	// b.StopTimer()

	// b.StartTimer()

	for i := 0; i < b.N; i++ {
		var allContacts []Contact
		err := Contactizer.GET(&allContacts, []int64{9})
		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}
