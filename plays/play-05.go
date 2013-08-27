// JSON Array parsing

package main

import (
	"encoding/json"
	"fmt"
)

const blob = `{"id":10,"name":"Invitation launch of GO 2011","description":"","created_at":"2012-04-05 12:47:15","updated_at":"2012-04-05 12:47:15","letter_template":null,"user_id":17,"share":1,"ids":[717,38,467,19,23,29,455,491,16,454,36,4,713,545,39,20,544,30,33,569,22,17,334,1229,1738,1537,211,590,529,488,600,225,1237,1541,1686,296,233,583,853,172,77,457,80,1533,56,183,1606,65,297,223,240,190,556,289,203,1622,218,399,495]}`

type Item struct {
	Ids []int64
}

func main() {
	var item Item
	json.Unmarshal([]byte(blob), &item)
	fmt.Printf("> %v \n", item)
}
