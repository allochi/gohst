package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type NHISRecord struct {
	// ID     int64
	HHX    int64
	FMX    int64
	FPX    int64
	SEX    int64
	BMI    int64
	SLEEP  int64
	EDUC   int64
	HEIGHT int64
	WEIGHT int64
}

func main() {

	file, err := os.OpenFile("/Users/allochi/Downloads/Udacity Sheets/nhisdata.txt", os.O_RDONLY, 0666)
	if err != nil {
		panic("")
	}

	// fileReader := bufio.NewReader(file)
	csvFile := csv.NewReader(file)

	header, err := csvFile.Read()
	for index, value := range header {
		header[index] = strings.ToUpper(value)
	}

	var line []string
	_type := reflect.TypeOf(NHISRecord{})
	for err != io.EOF {
		var record NHISRecord
		_record := reflect.New(_type)
		// _record := reflect.ValueOf(record)
		line, err = csvFile.Read()
		for index, _ := range header {
			_value, err := strconv.ParseInt(line[index], 10, 64)
			if err != nil {
				fmt.Errorf("format")
			}
			_record.Elem().FieldByIndex([]int{index})
			fmt.Println(_value)
		}
		record = _record.Elem().Interface().(NHISRecord)
		fmt.Printf("%#v \n", record)
	}

}
