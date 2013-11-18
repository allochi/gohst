package main

import (
	"allochi/gohst"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"strings"
	"time"
)

type Person struct {
	Id   int64
	Name string
}

func main() {

	imams := []Person{
		Person{1, "Mohammad"},
		Person{2, "Ali"},
		Person{3, "Hasan"},
		Person{4, "Husain"},
		Person{12, "Hujjah"},
	}

	e1 := gohst.Clause{"country_id:int", "IN", []int64{56, 79, 28, 90, 10}}
	e2 := gohst.Clause{"title", "LIKE", "Cheif %"}
	e3 := gohst.Clause{"colors:text", "IN", []string{"Red", "Green", "Blue"}}
	e4 := gohst.Clause{"days:datetime", "IN", []time.Time{time.Now(), time.Now().AddDate(-1, 0, 0)}}
	e5 := gohst.Clause{"Imam", "IN", imams}

	rc := &gohst.RequestChain{}
	// rc.Where(e1).And(e2).And(e3)
	rc.Where(e1).And(e2).And(e3).Or(e4).And(e5)

	request := []string{
		"(", e1.Bake(), "AND", e2.Bake(), ")",
		"OR",
		"(", e3.Bake(), "AND", e4.Bake(), ")",
	}

	fmt.Printf("%s\n\n", rc.Bake())
	fmt.Printf("%s\n\n", strings.Join(request, " "))

	rc2 := &gohst.RequestChain{}
	rc2.WhereGroup(e1).And(e2).AndGroup(e3).Or(e4).AndGroup(e5)
	fmt.Printf("%s\n\n", rc2.Bake())

	contactsRQ := &gohst.RequestChain{}
	contactsRQ.Where(gohst.Clause{"country_id:int", "=", 20}).And(gohst.Clause{"title_id:int", "=", 4}).And(gohst.Clause{"telephone", "=", ""})
	fmt.Printf("%s\n\n", contactsRQ.Bake())

	// fmt.Println(rc.Bake())
	// spew.Dump(rc)

}
