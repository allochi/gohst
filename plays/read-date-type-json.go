package main

import (
	// "allochi/gohst"
	. "allochi/gohst/plays/models"
	"encoding/json"
	"fmt"
	"time"
)

func main() {

	testContact := Contact{}
	testContact.ArchivedAt = time.Now()
	testData, err := json.Marshal(testContact)
	fmt.Printf("Test: %s\n\n", testData)

	data := `{"id":12,"title_id":4,"first_name":"Pascal","last_name":"Desbiens","job_title_id":8,"organization_id":5,"department_id":null,"address_01":"Avenue de l'Ariana 5","address_02":"","address_03":"","city_id":1397,"postal_code":"CH-1202","country_id":195,"telephone":"022 919 92 24","mobile":"","fax":"022 919 92 33","email_1":"pascal.desbiens@international.gc.ca","email_2":"","created_at":"2010-06-02 06:56:20","updated_at":"2012-07-28 10:37:46","role":"to","comments":"","address_text":"Avenue de l'Ariana 5\r\n","delta":0,"updated_by_user_id":1,"managed_by_user_id":null,"archived_by_user_id":17,"archived_at":"2012-07-28 00:00:00","archived_comment":"","focal_point":1,"isOrganization":false,"title":"Mr.","job_title":"Counsellor, Humanitarian Affairs","country":"Switzerland","city":"CH-1202 Genève","organization":"Permanent Mission of Canada","department":null,"emails":["pascal.desbiens@international.gc.ca",""],"is_organization":false,"is_focal_point":true,"categories":["Governments"],"interests":["India"],"sectors":["Donors","Governments"]}`

	contact := Contact{}
	err = json.Unmarshal([]byte(data), &contact)
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	}

	fmt.Printf("Object: %#V\n", contact)

}
