package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	. "allochi/tcolor"
	// "github.com/davecgh/go-spew/spew"
	"fmt"
	"log"
	"time"
)

var tc = TColor
var Contactizer gohst.DataStore

func init() {

	DSPostJson := gohst.NewPostJson("allochi_contactizer", "allochi", "")
	DSPostJson.CheckCollections = true
	DSPostJson.AutoCreateCollections = true

	Contactizer.Register("Contactizer", DSPostJson)
	Contactizer.Connect()

}

func main() {

	defer Contactizer.Disconnect()

	start := time.Now()
	number := 10000

	for i := 0; i < number; i++ {
		var bog Bog
		bog.Name = "Allochi"
		bog.Messages = []string{"This is the first bog", "If another bog is created then it will be bogbog", "I don''t know if bogs are OK, but later we will have complex object tested"}
		bog.Tags = []string{"For", "example, based", "on", "experience, a", "candidate", "for", "an", "employee", "benefits", "management", "position", "might", "use", "the", "following", "resume", "keywords: employee", "benefit", "plans, CEBS, health", "care", "benefits, benefit", "policy, FMLA.A", "customer", "service", "representative", "could", "include: customer", "service, customer", "tracking", "system, computer", "skills, order", "entry", "experience."}
		bog.Link = "http://jobsearch.about.com/od/resumewriting/ig/Sections-of-a-Resume-Examples/Resume-Keywords.htm"

		err := Contactizer.Put(&bog)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	duration := time.Since(start)
	fmt.Printf("time: %vs\n", duration.Seconds())
	fmt.Printf("object per second: %v o/s\n", float64(number)/duration.Seconds())
	// spew.Dump(bog)
}
