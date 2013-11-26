package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
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

	var salaries []Salary
	trx1, err := Contactizer.Begin("StartMe")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	trx2, err := Contactizer.Begin("StartYou")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	trx3, err := Contactizer.Begin("StartYou")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf("%v\n", salaries)
	fmt.Printf("%v\n", trx1)
	fmt.Printf("%v\n", trx2)
	fmt.Printf("%v\n", trx3)

}
