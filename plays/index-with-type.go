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

}

func main() {

	Contactizer, err := gohst.GetDataStore("Contactizer")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	Contactizer.Index(Contact{}, "FirstName")
	Contactizer.Index(Contact{}, "Categories")
	err = Contactizer.Index(Contact{}, "ArchivedAt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	err = Contactizer.Index(Contact{}, "CreatedAt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	booking := time.Now()
	duration := time.Since(booking)
	fmt.Printf("Duration: %d ns done!\n", duration.Nanoseconds())

}
