package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	source := `[{"name":"Ali"},{"name":"Suaad"},{"name":"Ahmed"}]`
	var result interface{}
	err := json.Unmarshal([]byte(source), &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", result)
	fmt.Printf("There are %d objects\n", len(result.([]interface{})))
	for index, value := range result.([]interface{}) {
		fmt.Printf("[%d] %s\n", index, value.(map[string]interface{})["name"])
	}

}
