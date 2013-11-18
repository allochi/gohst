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

	for i := 0; i < b.N; i++ {
		var allContacts []Contact
		// request := &gohst.RequestChain{}
		// request.Where(gohst.Clause{"Id", "IN", []int64{9}})
		// err := Contactizer.Get(&allContacts, request)
		err := Contactizer.GetById(&allContacts, []int64{9})
		if err != nil {
			b.Fatalf("gohst error: %s", err)
		}
	}

}
