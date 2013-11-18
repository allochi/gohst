// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
	"testing"
	"time"
)

func TestRequestIN(t *testing.T) {

	ids := []int64{4, 5, 6, 7, 8, 9}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	query := request.Bake(Contact{})

	expected := "WHERE id IN (4,5,6,7,8,9)"

	if query != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, query)
	}

}

func TestRequestDate(t *testing.T) {

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"ArchivedAt", ">", time.Now()})
	query := request.Bake(Contact{})

	expected := fmt.Sprintf("WHERE _date(data,'archived_at') > '%s'", time.Now().Format("2006-01-02 15:04:05"))

	if query != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, query)
	}

}

func TestRequestTextArray(t *testing.T) {

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Categories", "<@", "{Governments}"})
	query := request.Bake(Contact{})

	expected := "WHERE '{Governments}' <@ _array(data,'categories')"

	if query != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, query)
	}

}
