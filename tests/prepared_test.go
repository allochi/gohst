// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	// "allochi/inflect"
	"database/sql"
	// "encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	// "strings"
	"testing"
	"time"
)

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

	db, _ = sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

}

func TestPreparedSelect(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	var contacts []Contact

	Contactizer.Index(Contact{}, "categories")

	categories := []string{"Governments", "Donors"}

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"categories", "<@", gohst.IN(categories, true)})

	requestPrepared := &gohst.RequestChain{}
	requestPrepared.Where(gohst.Clause{"categories", "<@", "$1"})
	Contactizer.Prepare("SelectByCategory", Contact{}, requestPrepared)

	checkTimeDirect := time.Now()
	Contactizer.Get(&contacts, request)
	durationDirect := time.Since(checkTimeDirect)
	fmt.Printf("[Direct  ]: %s\n", durationDirect)

	checkTimePrepared := time.Now()
	Contactizer.ExecutePrepared("SelectByCategory", &contacts, gohst.IN(categories, true))
	durationPrepared := time.Since(checkTimePrepared)
	fmt.Printf("[Prepared]: %s\n", durationPrepared)

	ratio := float64(durationPrepared) / float64(durationDirect)
	fmt.Printf("[Ratio   ]: %.3f\n", ratio)

	if ratio > 1 {
		t.Errorf("Expected prepared/direct time %f to be less than 1", ratio)
	}

}

func TestPreparedStatementsWithCasting(t *testing.T) {

	var expected int
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	db.QueryRow("select count(*) from json_contacts where (data->>'country_id')::int > 20;").Scan(&expected)

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	var contacts []Contact

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"country_id::int", ">", "$1"})
	Contactizer.Prepare("SelectByCountryId", Contact{}, request)

	err = Contactizer.ExecutePrepared("SelectByCountryId", &contacts, 20)
	got := len(contacts)
	if err != nil || got != expected {
		t.Errorf("Expected %d contacts got %d", expected, got)
	}

}
