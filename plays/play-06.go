// JSON Array parsing

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

func main() {

	ids := []int64{1, 2, 3, 4}
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.FormatInt(id, 10)
	}

	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	defer db.Close()

	// rows, err := db.Query("select id from json_contacts;")
	sqlStatement := fmt.Sprintf("select id from json_contacts where id in (%s);", strings.Join(idsStr, ","))
	rows, err := db.Query(sqlStatement)
	defer rows.Close()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		fmt.Printf("%d\n", id)
	}

}
