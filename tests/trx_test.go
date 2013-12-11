// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	_ "github.com/lib/pq"
	"testing"
)

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	postgres.CheckCollections = true
	postgres.AutoCreateCollections = true
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

}

func TestTrxHasUniqueNames(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	var err error

	_, err = Contactizer.Begin("A")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = Contactizer.Begin("B")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = Contactizer.Begin("B")
	if err == nil {
		t.Errorf("Expected an error transactions can't have the same name")
	}

}

func TestTrxScopeAndRollback(t *testing.T) {

	Contactizer, err := gohst.GetDataStore("Contactizer")
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Empty json_salaries
	err = Contactizer.Drop(Salary{}, true)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	trx, err := Contactizer.Begin("") // name will be timestamp
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	salary := Salary{}
	salary.Name = "Allochi"
	salary.Amount = 120
	Contactizer.Put__(trx, salary)

	// -- Get outside Trx
	var salaries []Salary
	err = Contactizer.Get(&salaries, []int64{})
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if len(salaries) > 0 {
		t.Errorf("Expected 0 objects but got %d\n", len(salaries))
	}

	// -- Get inside Trx
	salaries = []Salary{}
	err = Contactizer.Get__(trx, &salaries, []int64{})
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if len(salaries) != 1 {
		t.Errorf("Expected 1 object but got %d\n", len(salaries))
	}

	Contactizer.Rollback(trx)

	// Get after rollback
	salaries = []Salary{}
	err = Contactizer.Get(&salaries, []int64{})
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if len(salaries) > 0 {
		t.Errorf("Expected 0 objects but got %d\n", len(salaries))
	}

}
