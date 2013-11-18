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
	// createIndexOnCountryId(Contactizer)

	// timer := time.Now()
	// rawContacts, err := Contactizer.ExecuteRaw("empty_names()")
	// rawContacts, err := Contactizer.GetRaw(Contact{}, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	// rawContacts, err := Contactizer.GetRaw(Contact{}, []int64{9, 1, 8, 2, 6, 3, 7, 5}, "")

	// Query in 0.034391297s
	// 988 row -> 28,721 Row/s
	// Contactizer.GetRaw(Contact{}, []int64{}, "(data->>'country_id')::int DESC")

	// Query in 0.014067022s
	// Contactizer.GetRaw(Contact{}, []int64{}, "")

	// var contacts []Contact
	// Contactizer.Execute(&contacts, "empty_names()")
	// Contactizer.Get(&contacts, []int64{9, 1, 8, 2, 6, 3, 7, 5}, "(data->>'country_id')::int DESC")

	// Query All in 0.156085775s (4.54x Raw)
	// 988 row -> 6,330 obj/s
	// Contactizer.Get(&contacts, []int64{}, "(data->>'country_id')::int DESC")

	// withGo(Contactizer)

	// raw(Contactizer)
	// objects(Contactizer)
	// preparedObjects(Contactizer)
	preparedArrayFind(Contactizer)

}

func raw(Contactizer gohst.DataStore) {
	fmt.Println("Running raw()")
	// == no index no prepared:
	// == Query in 0.038814292s, 971883 chars were retrieved (17")

	// == index no prepared:
	// == Query in 0.036462497s, 971883 chars were retrieved (17")

	// createIndexOnTitleId(Contactizer)

	timer := time.Now()
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"country_id:int", ">", 20})
	result, err := Contactizer.GetRaw(Contact{}, request)
	if err != nil {
		fmt.Println(err)
		return
	}
	duration := time.Since(timer).Nanoseconds()
	fmt.Printf("Query in %vs, %d chars were retrieved\n", float64(duration)/float64(1000000000), len(result))
}

func objects(Contactizer gohst.DataStore) {
	var contacts []Contact
	timer := time.Now()
	// Query in 0.146086243s using empty request to grap all 988 record
	// Query in 0.00228221s for 8 objects (3.5k objs aprox.)
	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"country_id:int", ">", 20})
	// request.Where(gohst.Clause{"Id", "IN", []int64{9, 1, 8, 2, 6, 3, 7, 5}})
	Contactizer.Get(&contacts, request)

	// Contactizer.Get(&contacts, []int64{}, "")
	duration := time.Since(timer).Nanoseconds()

	fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
	fmt.Printf("Query in %vs\n", float64(duration)/float64(1000000000))

}

func preparedObjects(Contactizer gohst.DataStore) {
	var contacts []Contact

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"country_id:int", ">", "$1"})
	Contactizer.Prepare("SelectByCountryId", Contact{}, request)

	timer := time.Now()
	// 891 Contacts in 0.12906126s
	err := Contactizer.ExecutePrepared("SelectByCountryId", &contacts, 20)
	if err != nil {
		fmt.Println(err)
	}
	duration := time.Since(timer).Nanoseconds()

	fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
	fmt.Printf("Query in %vs\n", float64(duration)/float64(1000000000))

}

func preparedArrayFind(Contactizer gohst.DataStore) {
	var contacts []Contact

	request := &gohst.RequestChain{}
	request.Where(gohst.Clause{"categories:text[]", "@>", "$1"})
	Contactizer.Prepare("SelectByCategory", Contact{}, request)

	timer := time.Now()
	// 891 Contacts in 0.12906126s
	err := Contactizer.ExecutePrepared("SelectByCategory", &contacts, "{Governments, Donors}")
	if err != nil {
		fmt.Println(err)
	}
	duration := time.Since(timer).Nanoseconds()

	fmt.Printf("%d Contacts in %vs\n", len(contacts), float64(duration)/float64(1000000000))
	fmt.Printf("Query in %vs\n", float64(duration)/float64(1000000000))

}

func createIndexOnCountryId(Contactizer gohst.DataStore) {
	err := Contactizer.Index(Contact{}, "country_id", "int")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func createIndexOnTitleId(Contactizer gohst.DataStore) {
	err := Contactizer.Index(Contact{}, "title_id", "int")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func goroutines(Contactizer gohst.DataStore) {
	doneCount := 0

	go func() {
		timer := time.Now()
		var contacts []Contact
		request := &gohst.RequestChain{}
		request.Where(gohst.Clause{"Id", "IN", []int64{9, 1, 8, 2, 6, 3, 7, 5}})
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
		request.Where(gohst.Clause{"country_id:int", ">", 20})
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
