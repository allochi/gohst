package main

import "fmt"

// import "reflect"
import "time"
import "allochi/gohst"

type Contact struct {
	Id        int64  `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string
	Addresses []Address   `json:"Addresses"`
	Emails    []string    `json:"emails"`
	Logs      []time.Time `json:"logs"`
	CreatedAt time.Time   `json:"-"`
	Sex       bool
	Salary    float32
	Color     uint32
}

type Address struct {
	Street  string
	Zipcode string
}

type Fields struct {
	OrgInfo map[string]string
	DstInfo map[string]string
}

func main() {

	// si := Fields{}
	// si.OrgInfo = map[string]string{}
	// si.DstInfo = map[string]string{}

	contact := Contact{}
	// _value := reflect.ValueOf(contact)
	// _type := _value.Type()

	// for i := 0; i < _value.NumField(); i++ {
	// 	_field := _type.Field(i)
	// 	fmt.Println(_field.Type.String())
	// 	_fieldType := _field.Type.String()
	// 	_name := _field.Name
	// 	_jsonName := _field.Tag.Get("json")

	// 	si.OrgInfo[_name] = _fieldType
	// 	if !(_jsonName == "-" || _jsonName == "") {
	// 		si.DstInfo[_jsonName] = _fieldType
	// 	}

	// }

	fmt.Printf("%s\n", gohst.TypeOf(contact, "FirstName"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "Addresses"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "Logs"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "CreatedAt"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "Sex"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "Salary"))
	fmt.Printf("%s\n", gohst.TypeOf(contact, "Color"))

}
