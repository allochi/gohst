// test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

func TestReadData(t *testing.T) {

	var Contactizer gohst.PostJsonDataStore
	Contactizer = gohst.NewPostJson("allochi_contactizer", "allochi", "")
	Contactizer.Connect()
	// defer Contactizer.Disconnect()
	gohst.Register("Contactizer", Contactizer)

	var contacts []Contact

	var err error
	var expected int
	var got int

	ids := []int64{4, 5, 6, 7, 8, 9}
	expected = len(ids)
	err = gohst.GET("Contactizer", &contacts, ids)

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
	err = gohst.GET("Contactizer", &allContacts, []int64{})
	got = len(allContacts)
	if got != expected {
		t.Errorf("Get All: Expected %d contacts got %d", expected, got)
	}
}

// 17: BenchmarkReadData	    5000	    347401 ns/op
func BenchmarkReadData(b *testing.B) {

	b.StopTimer()

	var Contactizer gohst.PostJsonDataStore
	Contactizer = gohst.NewPostJson("allochi_contactizer", "allochi", "")
	Contactizer.Connect()
	// defer Contactizer.Disconnect()
	gohst.Register("Contactizer", Contactizer)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		var allContacts []Contact
		err := gohst.GET("Contactizer", &allContacts, []int64{9})
		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}
