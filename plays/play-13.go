package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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

var header []string
var line []string
var skipped [][]string
var ObjectPool Pool

func main() {

	file, err := os.OpenFile("/Users/allochi/Downloads/Udacity Sheets/nhisdata.csv", os.O_RDONLY, 0666)
	if err != nil {
		panic("")
	}

	// fileReader := bufio.NewReader(file)
	csvFile := csv.NewReader(file)

	header, err = csvFile.Read()
	for index, value := range header {
		header[index] = strings.ToUpper(value)
	}

	startTime := time.Now()
	for err != io.EOF {
		line, err = csvFile.Read()

		// Without Pool
		// withoutPool(line)

		// With Pool
		ObjectPool.Create(reflect.TypeOf(NHISRecord{}), 10, true)
		withPool(line)
	}
	duration := time.Since(startTime)

	fmt.Printf("It took: %fs\n", duration.Seconds())
	fmt.Printf("%v Skipped!\n%v\n", len(skipped), skipped)

}

func withoutPool(line []string) {
	var record NHISRecord
	_record := reflect.New(reflect.TypeOf(record))

	if len(line) < len(header) {
		skipped = append(skipped, line)
	} else {
		for index, _ := range header {
			_value, err := strconv.ParseInt(line[index], 10, 64)
			if err != nil {
				fmt.Errorf("format error")
			}
			_record.Elem().FieldByIndex([]int{index}).SetInt(_value)
		}
		record = _record.Elem().Interface().(NHISRecord)
		fmt.Printf("%#v \n", record)
	}
}

func withPool(line []string) {
	fmt.Printf("Pool(%d/%d)\n", ObjectPool.Next(), ObjectPool.Capacity())
	if len(line) < len(header) {
		skipped = append(skipped, line)
	} else {
		var _elem reflect.Value
		_record, err := ObjectPool.Pull()
		if err != nil {
			panic(err)
		}
		_elem = _record
		for index, _ := range header {
			_value, err := strconv.ParseInt(line[index], 10, 64)
			if err != nil {
				fmt.Errorf("format error")
			}
			_elem.FieldByIndex([]int{index}).SetInt(_value)
		}
		// record := _elem.Interface().(NHISRecord)
		// fmt.Printf("%#v \n", record)
	}
}

type Pool struct {
	typed    reflect.Type
	slice    reflect.Value
	capacity int
	next     int
	revive   bool
	dead     bool
}

func (pool *Pool) Create(typed reflect.Type, capacity int, revive bool) {

	pool.typed = typed
	pool.capacity = capacity
	pool.revive = revive
	pool.slice = reflect.MakeSlice(reflect.SliceOf(typed), capacity, capacity)

}

func (pool *Pool) Pull() (object reflect.Value, err error) {
	// fmt.Printf("Revived Cap:%d Next:%d Dead:%v\n", pool.capacity, pool.next, pool.dead)
	fmt.Printf("Pool Pull %4d of %4d (D:%v - C:%v): ", pool.next, pool.capacity, pool.dead, (!pool.dead && pool.next <= pool.slice.Cap()))
	switch {
	case !pool.dead && pool.next < pool.capacity:
		object = pool.slice.Index(pool.next)
		pool.next++
		if pool.next >= pool.capacity {
			pool.dead = true
		}
	case pool.dead && pool.revive:
		object = reflect.Indirect(reflect.New(pool.typed))
		pool.slice = reflect.Append(pool.slice, object)
		fmt.Printf("\n%v\n", reflect.Indirect(pool.slice).Interface())
		// pool.capacity = pool.slice.Cap()
		pool.dead = false
		// pool.next++
		fmt.Printf("Revived Cap:%d Next:%d Dead:%v\n", pool.slice.Cap(), pool.next, pool.dead)
	default:
		err = fmt.Errorf("Pool is dead and can't be revived, it has reached it's capacity of %d", pool.capacity)
	}

	return
}

func (pool *Pool) Dead() bool {
	return pool.dead
}

func (pool *Pool) Capacity() int {
	return pool.capacity
	// return pool.slice.Len()
}

func (pool *Pool) Next() int {
	return pool.next
}
