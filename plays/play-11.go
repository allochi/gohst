package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	"github.com/davecgh/go-spew/spew"
	"log"
)

var tc = TColor

var Contactizer gohst.DataStore

func init() {

	ContactizerJson := gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	ContactizerJson.CheckCollections = true
	ContactizerJson.AutoCreateCollections = true

	Contactizer.Register("Contactizer", ContactizerJson)
	Contactizer.Connect()
}

func main() {

	defer Contactizer.Disconnect()

	var bogs []Bog
	for i := 0; i < 8; i++ {
		var bog Bog
		Contactizer.Put(&bog)
		bogs = append(bogs, bog)
	}

	// First 4 bogs
	var firstIds []int64
	for i := 0; i < 8; i++ {
		firstIds = append(firstIds, bogs[i].Id)
	}

	// Delete without return deleted objects
	var bogs2Delete []Bog
	err := Contactizer.Delete(bogs2Delete, firstIds)
	// err := Contactizer.Delete(&bogs2Delete, firstIds)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Delete retrieved %d bogs", len(bogs2Delete))

	// Try to get them
	var deletedBogs []Bog
	Contactizer.Get(&deletedBogs, firstIds)
	log.Printf("Found %d of deleted bogs", len(deletedBogs))

	spew.Dump(bogs2Delete)
}
