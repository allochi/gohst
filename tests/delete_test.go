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
	postgres := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	postgres.CheckCollections = true
	postgres.AutoCreateCollections = true
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

}

func insertGreeksAlphabet() {
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	for _, greek := range greekAlphabet() {
		Contactizer.Put(greek)
	}
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
