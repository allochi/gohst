package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
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

	Contactizer.Drop(Salary{}, true)

	trx, err := Contactizer.Begin("StartMe")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	salary := Salary{}
	salary.Name = "Allochi"
	salary.Amount = 23000
	Contactizer.Put__(salary, trx)

	fmt.Println("-- Not with Trx")
	var salaries []Salary
	err = Contactizer.Get(&salaries, []int64{})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	for _, salary := range salaries {
		fmt.Printf("%v\n", salary)
	}

	fmt.Println("-- With Trx")
	salaries = []Salary{}
	err = Contactizer.Get__(&salaries, []int64{}, trx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	for _, salary := range salaries {
		fmt.Printf("%v\n", salary)
	}

	Contactizer.Commit(trx)

}
