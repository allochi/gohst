package main

import (
	// "allochi/inflect"
	// "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	// "reflect"
	// "strings"
	. "allochi/gohst/plays/models"
	"time"
)

type DataRow struct {
	Id         int64
	Data       []byte
	ArchivedAt time.Time
}

func main() {

	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	rows, err := db.Query("SELECT id, data, (data->>'archived_at')::TIMESTAMP from json_contacts;")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for rows.Next() {
		dataRow := DataRow{}
		rows.Scan(&dataRow.Id, &dataRow.Data, &dataRow.ArchivedAt)
		contact := Contact{}
		err = json.Unmarshal(dataRow.Data, &contact)
		// if err != nil {
		// 	fmt.Printf("Error: %s\n", err)
		// }
		contact.ArchivedAt = dataRow.ArchivedAt
		data, err := json.Marshal(contact)
		_, err = db.Exec("UPDATE json_contacts SET data = $1 WHERE id = $2", data, dataRow.Id)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}

}
