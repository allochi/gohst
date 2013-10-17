package main

import (
	"fmt"
	"os"
)

func main() {

	file, err := os.OpenFile("/Users/allochi/Downloads/Udacity Sheets/nhisdata.txt", os.O_RDONLY, 0666)
	fmt.Println(file)

}
