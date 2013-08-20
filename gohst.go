package gohst

var DataStore DataStoreDeligate

type DataStoreDeligate interface {
	PUT(interface{}) Response
	GET(interface{}, interface{}) Response
}

type Response struct {
	Message string
	Error   error
	Size    int
}

type Request struct {
	Ids    []int64
	Source string
}

func PUT(object interface{}) Response {
	return DataStore.PUT(object)
}

func GET(object interface{}, query interface{}) Response {
	return DataStore.GET(object, query)
}
