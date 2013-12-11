// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"allochi/inflect"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()
	// defer Contactizer.Disconnect()
}

func TestReadData(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	var contacts []Contact

	var err error
	var expected int
	var got int

	ids := []int64{4, 5, 6, 7, 8, 9}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	expected = len(ids)
	err = Contactizer.Get(&contacts, request)

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
	err = Contactizer.Get(&allContacts, &gohst.RequestChain{}) // empty request get all
	got = len(allContacts)
	if got != expected {
		t.Errorf("Get All: Expected %d contacts got %d", expected, got)
	}
}

func TestPreparedFindInArray(t *testing.T) {

	var expected int
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	db.QueryRow("SELECT count(*) FROM json_contacts WHERE '{Governments, Donors}' <@ _array(data,'categories');").Scan(&expected)

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	var contacts []Contact

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"categories", "@>", "$1"})
	Contactizer.Prepare("SelectByCategory", Contact{}, request)

	err = Contactizer.ExecutePrepared("SelectByCategory", &contacts, "{Governments, Donors}")
	got := len(contacts)
	if err != nil || got != expected {
		t.Errorf("Expected %d contacts got %d", expected, got)
	}

}

func TestJsonName(t *testing.T) {
	field := "FirstName"
	object := Contact{}
	gohst.ScanFields(object)
	expected := "first_name"
	got := gohst.JsonName(object, field)
	if got != expected {
		t.Errorf("Expected %s instead got %s", expected, got)
	}
}

// go test -run CreateIndex$ -v
func TestCreateIndex(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	object := Contact{}
	gohst.ScanFields(object)
	field := "FirstName"
	_name, _, _ := gohst.TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	_json := gohst.JsonName(object, field)
	if _json != "" {
		field = _json
	}

	indexName := fmt.Sprintf("%s_%s_idx", tableName, field)

	// Drop index if it exist
	db.Exec(fmt.Sprintf("DROP INDEX %s;", indexName))

	err = Contactizer.Index(object, field)
	if err != nil {
		t.Errorf("Couldn't create index: ", err)
	}

	var got int
	db.QueryRow("SELECT count(*) FROM pg_indexes WHERE tablename = $1 AND indexname = $2;", tableName, indexName).Scan(&got)
	if got != 1 {
		t.Errorf("Expected %s to exist, instead got %d rows in pg_indexes", indexName, got)
	}

}

func TestCreateIndexForDate(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	object := Contact{}
	gohst.ScanFields(object)
	field := "ArchivedAt"
	_name, _, _ := gohst.TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	_json := gohst.JsonName(object, field)
	if _json != "" {
		field = _json
	}

	indexName := fmt.Sprintf("%s_%s_idx", tableName, field)

	// Drop index if it exist
	db.Exec(fmt.Sprintf("DROP INDEX %s;", indexName))

	statement, _ := db.Prepare("EXPLAIN (ANALYZE true, COSTS true, FORMAT json) SELECT _date(data,'archived_at') FROM json_contacts WHERE _date(data,'archived_at') BETWEEN TIMESTAMP '2011-05-01 00:00:00' AND TIMESTAMP '2011-06-01 00:00:00';")

	durationB4Index := statementTotalRuntime(statement)
	fmt.Printf("[- index results]: %f ms\n", durationB4Index)

	err = Contactizer.Index(object, field)
	if err != nil {
		t.Errorf("Couldn't create index: ", err)
	}

	durationIndex := statementTotalRuntime(statement)
	fmt.Printf("[+ index results]: %f ms\n", durationIndex)

	if durationB4Index <= durationIndex {
		t.Errorf("Expected duration before index %s to be less than before an index %s", durationB4Index, durationIndex)
	}

	request := &gohst.RequestChain{}
	// request.Where(gohst.Clause{"archived_at", "<", "2011-05-01 00:00:00"})
	date, _ := time.Parse("2006-01-02 15:04:05", "2011-05-01 00:00:00")
	request.Where(gohst.Clause{"archived_at", "<", date})
	listTwo := []Contact{}
	Contactizer.Get(&listTwo, request)

	testTime, _ := time.Parse("2006-01-02 03:04:05", "2011-05-01 00:00:00")
	for _, contact := range listTwo {
		if contact.ArchivedAt.After(testTime) {
			t.Errorf("Expected all dates to be before (%s) and got (%s)", testTime, contact.ArchivedAt)
		}
	}

}

func TestCreateIndexForArray(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	object := Contact{}
	gohst.ScanFields(object)
	field := "Categories"
	_name, _, _ := gohst.TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	_json := gohst.JsonName(object, field)
	if _json != "" {
		field = _json
	}

	indexName := fmt.Sprintf("_array_%s_%s_idx", tableName, field)

	// DROP INDEX _array_json_contacts_categories_idx;
	// CREATE INDEX _array_json_contacts_categories_idx ON json_contacts USING GIN (_array(data,'categories'));
	sql := fmt.Sprintf("DROP INDEX %s;", indexName)
	db.Exec(sql)

	statement, _ := db.Prepare("EXPLAIN (ANALYZE true, COSTS true, FORMAT json) SELECT * FROM json_contacts WHERE '{Governments, Donors}' <@ _array(data,'categories');")

	durationB4Index := statementTotalRuntime(statement)
	fmt.Printf("[- index results]: %f ms\n", durationB4Index)

	err = Contactizer.Index(object, field)
	if err != nil {
		t.Errorf("Couldn't create index: ", err)
	}

	durationIndex := statementTotalRuntime(statement)
	fmt.Printf("[+ index results]: %f ms\n", durationIndex)

	if durationB4Index <= durationIndex {
		t.Errorf("Expected duration before index %s to be less than before an index %s", durationB4Index, durationIndex)
	}

}

func statementTotalRuntime(statement *sql.Stmt) float64 {
	var results []byte
	statement.QueryRow().Scan(&results)
	return totalRuntime(results)
}

func totalRuntime(results []byte) float64 {
	var data interface{}
	json.Unmarshal(results, &data)
	return data.([]interface{})[0].(map[string]interface{})["Total Runtime"].(float64)
}

func totalCost(results []byte) float64 {
	var data interface{}
	json.Unmarshal(results, &data)
	return data.([]interface{})[0].(map[string]interface{})["Plan"].(map[string]interface{})["Total Cost"].(float64)
}

func _TestInsertAndDelete(t *testing.T) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	bog := Bog{}
	bog.Name = "Ali Anwar"
	bog.Messages = []string{"Go is cool", "PostgreSQL is cool"}
	bog.Tags = []string{"Go", "PostgreSQL"}
	bog.Link = "http://www.allochi.com"

	// Return and ID
	err := Contactizer.Put(&bog)

	if err != nil {
		t.Fatalf("gohst error: %s", err)
	}

	if bog.Id <= 0 {
		t.Errorf("Put: Expected an ID > 0 got %d", bog.Id)
	}

	ids := []int64{bog.Id}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	err = Contactizer.Delete([]Bog{}, request)

	if err != nil {
		t.Fatalf("gohst error: %s", err)
	}

	var bogs []Bog
	request = &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	err = Contactizer.Get(&bogs, request)

	if len(bogs) > 0 {
		t.Errorf("Delete: Didn't expect an object got %d", len(bogs))
	}

}

func TestPreparedInsertAndDelete(t *testing.T) {
	t.Skip("Need to implement prepared INSERT & DELETE statements")
}

func TestAddAuxFunctions(t *testing.T) {
	// _array() needs to be typed
	// Also joint index and full text search
	t.Skip("Need to implement adding all Aux. Functions like _array(data,'categories')")
}
