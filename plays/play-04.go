package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "reflect"
	"log"
	"strings"
)

var tc = TColor

func init() {
	// Contactizer = gohst.PostJsonDataStore{"allochi_contactizer", "allochi", ""}
	// Contactizer = gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	// Contactizer.Connect()
	// Contactizer.CheckCollections = true
	// gohst.Register("Contactizer", Contactizer)

	// var Contactizer gohst.DataStore
	ContactizerJson := gohst.NewPostJson("dbname=allochi_contactizer user=allochi sslmode=disable")
	Contactizer.Register("Contactizer", ContactizerJson)
	Contactizer.Connect()
	// Contactizer.Container.CheckCollections = true
}

func main() {

	ds, _ := gohst.GetDataStore("Contactizer")
	fmt.Printf(">> ds <<\n%#v \n\n", ds)
	fmt.Printf(">> Contactizer <<\n%#v \n\n", Contactizer)
	defer ds.Disconnect()

	// MailingListsList()
	// SamplePostJSon()
	// ContactsOfMailingList()
	// InsertAContact()
	// AllMailingLists()
	AllContacts()
}

// var Contactizer gohst.PostJsonDataStore
var Contactizer gohst.DataStore

func InsertAContact() {

	var contact Contact
	contact.FirstName = "Allochi"
	contact.LastName = "AlMuwali"
	// Clean temp insertions
	// delete from json_contacts where id >= 2107;
	// select setval('json_contacts_id_seq',2107);
	// select * from json_contacts where id >= 2107;
	err := Contactizer.Put(contact)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func UpdateAContact() {
	var contacts []Contact
	Contactizer.Get(&contacts, []int64{1})
	contact := contacts[0]
	// spew.Dump(contact)
	Contactizer.Put(contact)
}

func ContactsOfMailingList() {

	// Get the mailing list
	var mailingLists []MailingList
	err := Contactizer.Get(&mailingLists, []int64{14})

	if err != nil {
		log.Fatalf("Error getting mailing list: %s", err)
	}

	// Get the contacts
	var contacts []Contact
	Contactizer.Get(&contacts, mailingLists[0].ContactIds)

	for _, contact := range contacts {
		InterestOfContact(contact)
	}

	fmt.Println(len(mailingLists), " mailingLists were retrieved")
	// fmt.Println(len(contacts), " contacts were retrieved")

}

func MailingListsList() {
	var mailingLists []MailingList
	recordIDs := []int64{1, 2, 3}
	Contactizer.Get(&mailingLists, recordIDs)
	for _, mailingList := range mailingLists {
		fmt.Printf("> %s (%d)\n", mailingList.Name, len(mailingList.ContactIds))
	}
}

func AllMailingLists() {
	var mailingLists []MailingList
	// err := Contactizer.Get(&mailingLists, []int64{4, 5, 6, 7, 8, 9})
	err := Contactizer.Get(&mailingLists, []int64{})
	if err != nil {
		log.Fatalf("... %s", err)
	}
	for _, mailingList := range mailingLists {
		fmt.Printf("[%2d] %s \n", mailingList.Id, mailingList.Name)
	}
}

func AllContacts() {
	var contacts []Contact
	// err := Contactizer.Get(&contacts, []int64{4, 5, 6, 7, 8, 9})
	// err := Contactizer.Get(&contacts, []int64{})
	ds, err := gohst.GetDataStore("Contactizer")
	if err != nil {
		fmt.Errorf("err: %s", err)
	}
	err = ds.Get(&contacts, []int64{})
	if err != nil {
		log.Fatalf("... %s", err)
	}
	fmt.Printf("There are %d contacts \n\n", len(contacts))
	fmt.Printf("Contactizer = %#v \n\n", ds)
	// for _, contact := range contacts {
	// 	// fmt.Printf("[%4d] %s \n", contact.Id, contact.Name())
	// 	// fmt.Printf("%#v\n\n", contact)
	// 	// fmt.Printf("%+v\n\n", contact)
	// }
}

func SamplePostJSon() {

	var contacts []Contact
	Contactizer.Get(&contacts, []int64{1, 2, 3, 4})
	for _, contact := range contacts {
		InterestOfContact(contact)
	}

}

// func organizations(contact Contact) {
// 	if contact.IsOrganization {
// 		fmt.Printf("> %s\n %-18s %-24s\n\n", contact.Name(), contact.Country, contact.City)
// 	}
// }

// func categorizedContact(contact Contact) {
// 	if len(contact.Categories) > 0 {
// 		fmt.Printf("%-64s [ %s ]\n", contact.Name(), strings.Join(contact.Categories, " | "))
// 	} else {
// 		fmt.Printf("%s\n", contact.Name())
// 	}
// }

// func sectoredContact(contact Contact) {
// 	if len(contact.Sectors) > 0 {
// 		fmt.Printf("%-64s [ %s ]\n", contact.Name(), strings.Join(contact.Sectors, " | "))
// 	} else {
// 		fmt.Printf("%s\n", contact.Name())
// 	}
// }

func InterestOfContact(contact Contact) {
	if len(contact.Interests) > 0 {
		fmt.Printf("%-64s [ %s ]\n", contact.Name(), strings.Join(contact.Interests, " | "))
	} else {
		fmt.Printf("%s\n", contact.Name())
	}
}

// func SampleArray() {
// 	allochi := Person{"Ali", "Anwar", 40, "Switzerland"}
// 	vanessa := Person{"Vanessa", "-", 26, "Switzerland"}

// 	// people = append(people, vanessa)

// 	gohst.DataStore = People
// 	fmt.Printf("Main (people b4Put): %s\n", People)
// 	response := gohst.Put(allochi)
// 	fmt.Printf("%s\n", response.Message)
// 	if response.Message == "OK" {
// 		fmt.Printf("$v inserted", allochi)
// 	}
// 	if gohst.Put(vanessa).Message == "OK" {
// 		fmt.Printf("$v inserted", vanessa)
// 	}
// 	fmt.Printf("Main (people A8Put): %s\n", People)

// 	gohst.DataStore = others
// 	fmt.Printf("Main (others b4Put): %s\n", others)
// 	gohst.Put(allochi)
// 	fmt.Printf("Main (others A8Put): %s\n", others)
// }

// // --------------------------------------------------------------------------------
// // Array Data Stores
// // --------------------------------------------------------------------------------
// type PeopleDataStore []Person

// var People PeopleDataStore

// func (ds PeopleDataStore) Put(object interface{}) (response gohst.Response) {
// 	_object := reflect.ValueOf(object)
// 	_people := reflect.Indirect(reflect.ValueOf(&People))
// 	_people.Set(reflect.Append(_people, _object))

// 	response.Message = "Ok"
// 	response.Error = nil
// 	response.Size = 0
// 	return
// }

// func (ds PeopleDataStore) Get(interface{}, interface{}) (response gohst.Response) {
// 	return
// }

// type OthersDataStore []Person

// var others OthersDataStore

// func (ds OthersDataStore) Put(object interface{}) (response gohst.Response) {
// 	_object := reflect.ValueOf(object)
// 	_others := reflect.Indirect(reflect.ValueOf(&others))
// 	_others.Set(reflect.Append(_others, _object))

// 	response.Message = "Ok"
// 	response.Error = nil
// 	response.Size = 0
// 	return
// }

// func (ds OthersDataStore) Get(interface{}, interface{}) (response gohst.Response) {
// 	return
// }
