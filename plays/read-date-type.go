package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
	// "time"
)

func init() {

	ContactizerJson := gohst.NewPostJson("allochi_contactizer", "allochi", "")
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

	// err = Contactizer.Index(Contact{}, "ArchivedAt")
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// }

	// checkTime := time.Now()
	// duration := time.Since(checkTime)

	var contacts []Contact
	ids := []int64{12, 15, 21}
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"Id", "IN", ids})
	err = Contactizer.Get(&contacts, request)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for _, contact := range contacts {
		fmt.Printf("%v\n", contact.ArchivedAt)
	}

}
