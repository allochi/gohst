package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"fmt"
	"time"
)

func init() {

	ContactizerJson := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	ContactizerJson.CheckCollections = true
	ContactizerJson.AutoCreateCollections = true
	gohst.Register("Contactizer", ContactizerJson)

}

func main() {

	Contactizer, err := gohst.GetDataStore("Contactizer")
	if err != nil {
		fmt.Println(err)
	}
	err = Contactizer.Connect()
	if err != nil {
		fmt.Println(err)
	}
	defer Contactizer.Disconnect()

	// Contactizer.Get(&contacts, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	err = Contactizer.Index(Contact{}, "country_id", "int")
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	timer := time.Now()
	// rawContacts, err := Contactizer.ExecuteRaw("empty_names()")
	// rawContacts, err := Contactizer.GetRaw(Contact{}, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	// rawContacts, err := Contactizer.GetRaw(Contact{}, []int64{9, 1, 8, 2, 6, 3, 7, 5}, "")

	// Query in 0.034391297s
	// 988 row -> 28,721 Row/s
	// Contactizer.GetRaw(Contact{}, []int64{}, "(data->>'country_id')::int DESC")

	// Query in 0.014067022s
	// Contactizer.GetRaw(Contact{}, []int64{}, "")

	var contacts []Contact
	// Contactizer.Execute(&contacts, "empty_names()")
	// Contactizer.Get(&contacts, []int64{9, 1, 8, 2, 6, 3, 7, 5}, "(data->>'country_id')::int DESC")

	// Query All in 0.156085775s (4.54x Raw)
	// 988 row -> 6,330 obj/s
	// Contactizer.Get(&contacts, []int64{}, "(data->>'country_id')::int DESC")

	// Query in 0.146086243s using empty request to grap all 988 record
	// Query in 0.00228221s for 8 objects (3.5k objs aprox.)
	request := &gohst.RequestChain{}
	request.Where(gohst.Entry{"country_id:int", ">", 20})
	// request.Where(gohst.Entry{"Id", "IN", []int64{9, 1, 8, 2, 6, 3, 7, 5}})
	Contactizer.Get(&contacts, request)

	// Contactizer.Get(&contacts, []int64{}, "")
	duration := time.Since(timer).Nanoseconds()

	fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
	for _, contact := range contacts {
		fmt.Printf("%d", contact.Id)
	}
	fmt.Println()

	// fmt.Printf("%s\n\n in %vs\n", rawContacts, float64(duration)/float64(1000000000))
	fmt.Printf("Query in %vs\n", float64(duration)/float64(1000000000))

	withGo(Contactizer)

}

func withGo(Contactizer gohst.DataStore) {
	doneCount := 0

	go func() {
		timer := time.Now()
		var contacts []Contact
		request := &gohst.RequestChain{}
		request.Where(gohst.Entry{"Id", "IN", []int64{9, 1, 8, 2, 6, 3, 7, 5}})
		Contactizer.Get(&contacts, request)
		duration := time.Since(timer).Nanoseconds()
		fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
		for _, contact := range contacts {
			fmt.Printf("%d, ", contact.Id)
		}
		fmt.Println()
		doneCount++
	}()

	go func() {
		timer := time.Now()
		var contacts []Contact
		request := &gohst.RequestChain{}
		request.Where(gohst.Entry{"country_id:int", ">", 20})
		Contactizer.Get(&contacts, request)
		duration := time.Since(timer).Nanoseconds()
		fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
		for _, contact := range contacts {
			fmt.Printf("%d, ", contact.Id)
		}
		fmt.Println()
		doneCount++
	}()

	// var anything string
	// fmt.Scanln(&anything)
	for doneCount < 2 {
		time.Sleep(10 * time.Nanosecond)
	}
}
