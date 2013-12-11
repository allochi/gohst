// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
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

func TestPutOne(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	greekAlphabet := greekAlphabet()

	var err error
	err = Contactizer.Put(greekAlphabet[0])
	err = Contactizer.Put(greekAlphabet[1])
	err = Contactizer.Put(&greekAlphabet[2])

	if greekAlphabet[2].Id == 0 {
		t.Errorf("Expected object to have an Id, instead got %d", greekAlphabet[2].Id)
	}

	var greeks []Greek
	ids := []int64{1, 2, 3}
	err = Contactizer.Get(&greeks, ids)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	expected := len(ids)
	got := len(greeks)
	if got != expected {
		t.Errorf("Expected %d greeks got %d instead", expected, got)
	}

	cleanup()

}

func TestPutSlice(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	greekAlphabet := greekAlphabet()

	var err error
	err = Contactizer.Put(greekAlphabet)

	var greeks []Greek
	err = Contactizer.GetAll(&greeks)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	expected := len(greekAlphabet)
	got := len(greeks)
	// fmt.Printf("[exp] %d, [got] %d\n", expected, got)
	if got != expected {
		t.Errorf("Expected %d greeks got %d instead", expected, got)
	}

	cleanup()

}

func TestPutPointers(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	greekAlphabet := greekAlphabet()

	var err error
	err = Contactizer.Put(&greekAlphabet)

	var greeks []Greek
	err = Contactizer.GetAll(&greeks)

	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	expected := len(greekAlphabet)
	got := len(greeks)
	// fmt.Printf("[exp] %d, [got] %d\n", expected, got)
	if got != expected {
		t.Errorf("Expected %d greeks got %d instead", expected, got)
	}

	for _, object := range greekAlphabet {
		fmt.Sprintf("%#v\n", object)
		if object.Id == 0 {
			t.Errorf("Expected object.Id to have a value")
		}
	}

	cleanup()

}

func TestPutUpdate(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	greekAlphabet := greekAlphabet()

	var err error
	err = Contactizer.Put(&greekAlphabet)
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	var greeks []Greek
	err = Contactizer.GetAll(&greeks)
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	for i, _ := range greeks {
		greeks[i].Name = "O...O"
	}

	err = Contactizer.Put(&greeks)
	if err != nil {
		t.Errorf("Database Error: %s", err)
	}

	var ugreeks []Greek
	err = Contactizer.GetAll(&ugreeks)

	expected := len(greekAlphabet)
	got := len(ugreeks)
	if got != expected {
		t.Errorf("Expected %d greeks got %d instead", expected, got)
	}

	for _, object := range greekAlphabet {
		fmt.Sprintf("%#v\n", object)
		if object.Id == 0 {
			t.Errorf("Expected object.Id to have a value")
		}
	}

	cleanup()

}
