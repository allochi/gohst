package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	"github.com/davecgh/go-spew/spew"
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

	var bogs []Bog
	for i := 0; i < 8; i++ {
		var bog Bog
		gohst.PUT("Contactizer", &bog)
		bogs = append(bogs, bog)
	}

	// First 4 bogs
	var firstIds []int64
	for i := 0; i < 8; i++ {
		firstIds = append(firstIds, bogs[i].Id)
	}

	// Delete without return deleted objects
	var bogs2Delete []Bog
	err := gohst.DELETE("Contactizer", bogs2Delete, firstIds)
	// err := gohst.DELETE("Contactizer", &bogs2Delete, firstIds)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Delete retrieved %d bogs", len(bogs2Delete))

	// Try to get them
	var deletedBogs []Bog
	gohst.GET("Contactizer", &deletedBogs, firstIds)
	log.Printf("Found %d of deleted bogs", len(deletedBogs))

	spew.Dump(bogs2Delete)
}
