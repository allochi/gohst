package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	"fmt"
	"log"
	"time"
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

	ids := []int64{}

	var bogs1 []Bog
	var done1 bool
	go func() {
		log.Println("In {Go1}")
		Contactizer.Get(&bogs1, ids)
		log.Printf("Go1: There were %d bogs", len(bogs1))
		done1 = true
	}()

	var bogs2 []Bog
	var done2 bool
	go func() {
		log.Println("In {Go2}")
		Contactizer.Get(&bogs2, ids)
		log.Printf("Go2: There were %d bogs", len(bogs2))
		done2 = true
	}()

	for !(done1 || done2) {
		fmt.Println("waiting...")
		time.Sleep(10 * time.Millisecond)
	}

	// var input string
	// fmt.Scanln(&input)
	// fmt.Println("done")

}
