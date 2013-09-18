package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	// "github.com/davecgh/go-spew/spew"
	"log"
)

var tc = TColor
var Contactizer gohst.PostJsonDataStore

func init() {
	Contactizer = gohst.NewPostJson("allochi_contactizer", "allochi", "")
	Contactizer.Connect()
	Contactizer.CheckCollections = true
	Contactizer.AutoCreateCollections = true
	gohst.Register("Contactizer", Contactizer)
}

func main() {

	defer Contactizer.Disconnect()

	for i := 0; i < 1000; i++ {
		var bog Bog
		bog.Name = "Allochi"
		bog.Messages = []string{"This is the first bog", "If another bog is created then it will be bogbog", "I don''t know if bogs are OK, but later we will have complex object tested"}

		err := gohst.PUT("Contactizer", &bog)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	// spew.Dump(bog)
}
