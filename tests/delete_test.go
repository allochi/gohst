// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

var db *sql.DB
var greekAlphabet []Greek

type Greek struct {
	Id        int64 `json:"-"`
	Name      string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	postgres.CheckCollections = true
	postgres.AutoCreateCollections = true
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

	db, _ = sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	greekAlphabet = []Greek{
		Greek{0, "Αα Alpha", time.Now(), time.Now()},
		Greek{0, "Ββ Beta", time.Now(), time.Now()},
		Greek{0, "Γγ Gamma", time.Now(), time.Now()},
		Greek{0, "Δδ Delta", time.Now(), time.Now()},
		Greek{0, "Εε Epsilon", time.Now(), time.Now()},
		Greek{0, "Ζζ Zeta", time.Now(), time.Now()},
		Greek{0, "Ηη Eta", time.Now(), time.Now()},
		Greek{0, "Θθ Theta", time.Now(), time.Now()},
		Greek{0, "Ιι Iota", time.Now(), time.Now()},
		Greek{0, "Κκ Kappa", time.Now(), time.Now()},
		Greek{0, "Λλ Lambda", time.Now(), time.Now()},
		Greek{0, "Μμ Mu", time.Now(), time.Now()},
		Greek{0, "Νν Nu", time.Now(), time.Now()},
		Greek{0, "Ξξ Xi", time.Now(), time.Now()},
		Greek{0, "Οο Omicron", time.Now(), time.Now()},
		Greek{0, "Ππ Pi", time.Now(), time.Now()},
		Greek{0, "Ρρ Rho", time.Now(), time.Now()},
		Greek{0, "Σσ Sigma", time.Now(), time.Now()},
		Greek{0, "Ττ Tau", time.Now(), time.Now()},
		Greek{0, "Υυ Upsilon", time.Now(), time.Now()},
		Greek{0, "Φφ Phi", time.Now(), time.Now()},
		Greek{0, "Χχ Chi", time.Now(), time.Now()},
		Greek{0, "Ψψ Psi", time.Now(), time.Now()},
		Greek{0, "Ωω Omega", time.Now(), time.Now()},
	}

}

func insertGreeksAlphabet() {
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	for _, greek := range greekAlphabet {
		Contactizer.Put(greek)
	}
}

func cleanup() {
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	Contactizer.Drop(Greek{}, true)
}

func TestDeleteOne(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	insertGreeksAlphabet()

	var greeks []Greek
	ids := []int64{1}
	err := Contactizer.Get(&greeks, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	err = Contactizer.Delete(greeks[0])
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	greeks = []Greek{}
	err = Contactizer.Get(&greeks, ids)

	// var expected int
	// db.QueryRow("select count(*) from json_greeks where id in (4, 5, 6, 7, 8, 9);").Scan(&expected)

	got := len(greeks)
	if got > 0 {
		t.Errorf("Expected no greeks got %d", got)
	}

	cleanup()

}

func TestDeleteSlice(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	insertGreeksAlphabet()

	var greeks []Greek
	ids := []int64{1, 2, 3, 4}
	err := Contactizer.Get(&greeks, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	err = Contactizer.Delete(greeks)
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	greeks = []Greek{}
	err = Contactizer.Get(&greeks, ids)

	got := len(greeks)
	if got > 0 {
		t.Errorf("Expected no greeks got %d", got)
	}

	cleanup()

}

func TestDeletePointers(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	insertGreeksAlphabet()

	var greeks []Greek
	ids := []int64{1, 2, 3, 4}
	err := Contactizer.Get(&greeks, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	err = Contactizer.Delete(&greeks)
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	greeks = []Greek{}
	err = Contactizer.Get(&greeks, ids)

	got := len(greeks)
	if got > 0 {
		t.Errorf("Expected no greeks got %d", got)
	}

	// One object pointer
	greeks = []Greek{}
	err = Contactizer.Get(&greeks, []int64{5})

	err = Contactizer.Delete(&greeks[0])
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	greeks = []Greek{}
	err = Contactizer.Get(&greeks, []int64{5})

	got = len(greeks)
	if got > 0 {
		t.Errorf("Expected no greeks got %d", got)
	}

	cleanup()

}
