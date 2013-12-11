package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
	"time"
)

func init() {

	ContactizerJson := gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	ContactizerJson.CheckCollections = true
	ContactizerJson.AutoCreateCollections = true
	gohst.Register("Contactizer", ContactizerJson)
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	Contactizer.Connect()

}

func main() {

	Contactizer, err := gohst.GetDataStore("Contactizer")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	numObjects := 3187

	checkTime := time.Now()

	for i := 0; i < numObjects; i++ {
		var contacts []Contact
		ids := []int64{9}
		request := &gohst.RequestChain{}
		request.Where(gohst.Clause{"Id", "IN", ids})
		err = Contactizer.Get(&contacts, request)
		// Contactizer.GetById(&contacts, ids)
		// Contactizer.GetRawById(Contact{}, ids)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}

	duration := time.Since(checkTime)
	fmt.Printf("[Duration for %d]: %s\n", numObjects, duration)
	// fmt.Printf("[Duration]: %s\n", duration)

}
