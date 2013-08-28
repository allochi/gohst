package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "reflect"
	"strings"
)

var tc = TColor

func init() {
	Contacts = gohst.PostJsonDataStore{"allochi_contactizer", "allochi", ""}
	MailingLists = gohst.PostJsonDataStore{"allochi_contactizer", "allochi", ""}
	gohst.Register("Contacts", Contacts)
	gohst.Register("MailingLists", MailingLists)
}

func main() {
	// MailingListsList()
	// SamplePostJSon()
	// ContactsOfMailingList()
	InsertAContact()
}

var Contacts gohst.PostJsonDataStore
var MailingLists gohst.PostJsonDataStore

func InsertAContact() {

	var contact Contact
	contact.FirstName = "Allochi"
	contact.LastName = "AlMuwali"
	gohst.PUT("Contacts", contact)
}

func UpdateAContact() {
	var contacts []Contact
	gohst.GET("Contacts", &contacts, []int64{1})
	contact := contacts[0]
	// spew.Dump(contact)
	gohst.PUT("Contacts", contact)
}

func ContactsOfMailingList() {

	// Get the mailing list
	var mailingLists []MailingList
	gohst.GET("MailingLists", &mailingLists, []int64{14})

	// Get the contacts
	var contacts []Contact
	gohst.GET("Contacts", &contacts, mailingLists[0].ContactIds)

	for _, contact := range contacts {
		InterestOfContact(contact)
	}

	fmt.Println(len(contacts), " contacts were retrieved")

}

func MailingListsList() {
	var mailingLists []MailingList
	recordIDs := []int64{1, 2, 3}
	gohst.GET("MailingLists", &mailingLists, recordIDs)
	for _, mailingList := range mailingLists {
		fmt.Printf("> %s (%d)\n", mailingList.Name, len(mailingList.ContactIds))
	}
}

func SamplePostJSon() {

	var contacts []Contact
	gohst.GET("Contacts", &contacts, []int64{1, 2, 3, 4})
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
// 	fmt.Printf("Main (people b4PUT): %s\n", People)
// 	response := gohst.PUT(allochi)
// 	fmt.Printf("%s\n", response.Message)
// 	if response.Message == "OK" {
// 		fmt.Printf("$v inserted", allochi)
// 	}
// 	if gohst.PUT(vanessa).Message == "OK" {
// 		fmt.Printf("$v inserted", vanessa)
// 	}
// 	fmt.Printf("Main (people A8PUT): %s\n", People)

// 	gohst.DataStore = others
// 	fmt.Printf("Main (others b4PUT): %s\n", others)
// 	gohst.PUT(allochi)
// 	fmt.Printf("Main (others A8PUT): %s\n", others)
// }

// // --------------------------------------------------------------------------------
// // Array Data Stores
// // --------------------------------------------------------------------------------
// type PeopleDataStore []Person

// var People PeopleDataStore

// func (ds PeopleDataStore) PUT(object interface{}) (response gohst.Response) {
// 	_object := reflect.ValueOf(object)
// 	_people := reflect.Indirect(reflect.ValueOf(&People))
// 	_people.Set(reflect.Append(_people, _object))

// 	response.Message = "Ok"
// 	response.Error = nil
// 	response.Size = 0
// 	return
// }

// func (ds PeopleDataStore) GET(interface{}, interface{}) (response gohst.Response) {
// 	return
// }

// type OthersDataStore []Person

// var others OthersDataStore

// func (ds OthersDataStore) PUT(object interface{}) (response gohst.Response) {
// 	_object := reflect.ValueOf(object)
// 	_others := reflect.Indirect(reflect.ValueOf(&others))
// 	_others.Set(reflect.Append(_others, _object))

// 	response.Message = "Ok"
// 	response.Error = nil
// 	response.Size = 0
// 	return
// }

// func (ds OthersDataStore) GET(interface{}, interface{}) (response gohst.Response) {
// 	return
// }
