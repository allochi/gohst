// go test -bench=".*" ./tests
package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"testing"
)

func init() {

	var Contactizer gohst.DataStore
	postgres := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	gohst.Register("Contactizer", postgres)
	Contactizer, _ = gohst.GetDataStore("Contactizer")
	Contactizer.Connect()
	// defer Contactizer.Disconnect()
}

// ==========
// 17: BenchmarkReadData	    5000	    347401 ns/op (indirect access)
// 17: BenchmarkReadData	    5000	    315161 ns/op (3172 op/s - direct access)
// 17: BenchmarkReadData	    5000	    314334 ns/op (3181 op/s - direct access)
// ==========
// 13: BenchmarkReadData	    5000	    348989 ns/op (indirect access)
// 13: BenchmarkReadData	    5000	    361726 ns/op (indirect access)
// 13: BenchmarkReadData	   10000	    342638 ns/op (2918 op/s - direct access)
// 13: BenchmarkReadData	   10000	    313696 ns/op (3187 op/s - direct access)
// ==========
func BenchmarkReadData(b *testing.B) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	// request := &gohst.RequestChain{}
	// request.Where(gohst.Clause{"Id", "=", []int64{9}})
	// request.Where(gohst.Clause{"Id", "=", "$1"})
	// Contactizer.Prepare("One", Contact{}, request)

	for i := 0; i < b.N; i++ {
		var allContacts []Contact
		// BenchmarkReadData	    5000	    386074 ns/op
		// Contactizer.Get(&allContacts, request)

		// BenchmarkReadData	    5000	    266317 ns/op (3755 objects)
		// BenchmarkReadData	   10000	    268859 ns/op (17")
		// BenchmarkReadData	    5000	    273319 ns/op (13")
		// BenchmarkReadData	    5000	    311151 ns/op :( on 13"
		err := Contactizer.Get(&allContacts, []int64{9})

		// BenchmarkReadData	   10000	    281835 ns/op
		// Contactizer.ExecutePrepared("One", &allContacts, 9)

		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}

func BenchmarkGetAll(b *testing.B) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")

	for i := 0; i < b.N; i++ {
		var allContacts []Contact
		// BenchmarkGetAll	      20	 118348888 ns/op (17")
		// BenchmarkGetAll	      20	 116377441 ns/op (17")
		// BenchmarkGetAll	      20	  78797654 ns/op (13")
		// BenchmarkGetAll	      20	  95101618 ns/op :( on 13"
		err := Contactizer.Get(&allContacts)
		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}

func BenchmarkGetByArray(b *testing.B) {

	Contactizer, _ := gohst.GetDataStore("Contactizer")
	Contactizer.Index(Contact{}, "categories")
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"categories", "@>", "$1"})
	Contactizer.Prepare("SelectByCategory", Contact{}, request)

	for i := 0; i < b.N; i++ {
		var contacts []Contact
		// BenchmarkGetByArray	      50	  29549293 ns/op (17", about 338 times a second)
		// BenchmarkGetByArray	     500	   4235805 ns/op (17")
		// BenchmarkGetByArray	     500	   2935591 ns/op (13")
		// BenchmarkGetByArray	     500	   3406299 ns/op  :( on 13"
		err := Contactizer.ExecutePrepared("SelectByCategory", &contacts, "{Governments, Donors}")
		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}
