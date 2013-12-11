# gohst

gohst started as an ORM project, and then turned to a framework to use PostgreSQL json features to store objects, something that can have open schema, and support transaction, looking into the NoSQL world I couldn't find something that fits the way that I want things to work, so gohst turned from ambitious attempt to a practical tool.

Currently gohst works only with PostgreSQL, the intension was to support different databases, but since a decision was made to use PostgreSQL json extension, most of the code currently assumes that it's PostgreSQL that gohst is talking to, in the future this may change, and the interface may support this.

gohst also has been benchmarked on different laptops, on my MacBook Air, I'm able to pull avg. 3500 object/sec, medium sized, it depends on reflect and encoding/json packages, now with the release of GO 1.2 the json part should be faster.

## Behind the scenes

`gohst` converts a Go object to json and warps it with a predefined record object, then stores that in a table. 

The table names are the plural form of the struct name prefixed with `json_`, so an object of type `Contact` will be saved in the table `json_contacts`

```
type Contact struct {
	Id             int64     `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
```

As Go json package ignores fields with `json:"-"`, for raw queries, when the query result is just a pure json string, which can be forwarded directly to another server or a browser, you may need the `Id` field, or even the dates for sorting, therefor don't use `json:"-"`. 

**Structs to be used with **`gohst`** should always have these three fields.**

## About fields

Adding a new field to an object is nothing more than adding it to the struct, so to complete our Contact struct we can add more fields like so

```
type Contact struct {
	Id             int64     `json:"-"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Title          string    `json:"title"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	PostalCode     string    `json:"postal_code"`
	Interests      []string  `json:"interests"`
	ArchivedAt     time.Time `json:"archived_at"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
```
gohst should support all type of fields that can be converted by the Go json package, including array of primitive types.

Alternative json style fields names can be used using the `json` tag.

```
type Contact struct {
	Id             int64     `json:"-"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
```
translate to

```
{'first_name':'','last_name':''}
```

and

```
type Contact struct {
	Id             int64    
	FirstName      string   
	LastName       string   
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
```

translate to

```
{'Id':'', 'FirstName':'', 'LastName':'', 'CreatedAt':'', 'UpdatedAt':''}

```

## Removing fields

When a field is removed from the struct definition, it doesn't get removed automatically from the objects saved in the data store, to clean them up, objects have to be refreshed by pulling them from the database and saving them again.

## Connecting to database

gohst is using PostgreSQL, and it uses "github.com/lib/pq" as a driver, so all connection string parameters should be checked [here](http://godoc.org/github.com/lib/pq).

```
postgres := gohst.NewPostJson("dbname=contactizer user=allochi sslmode=disable")
postgres.CheckCollections = true
postgres.AutoCreateCollections = true
gohst.Register("Contactizer", postgres)
Contactizer, err := gohst.GetDataStore("Contactizer")
Contactizer.Connect()
```
With a new database and while development, it's better to set `CheckCollections` and `AutoCreateCollections` to `true`, this will check if the supporting tables exist and create them if not. 

gohst can drop tables with `gohst.Drop()`, if the there is a need to drop tables in production, these two options need to be `true`. 

As a note, `gohst.Drop()` doesn't always delete the table completely from the database, if the confirmation flag is not passed, then it rename the table by adding a timestamp of when `gohst.Drop()` was called.

`gohst.Register()` registers the database in gohst by name, so that it can be retrieved and used in different parts of the code, data stores names have to be unique. Multiple data stores can be used for different purposes, for example regional data partition.

```
postgres_eu := gohst.NewPostJson("dbname=contactizer_europe user=allochi sslmode=disable")
postgres_af := gohst.NewPostJson("dbname=contactizer_africa user=allochi sslmode=disable")
gohst.Register("Contactizer Europe", postgres_eu)
gohst.Register("Contactizer Africa", postgres_af)
```

Finally, a call to `gohst.Connect()` is needed

```
Contactizer, err := gohst.GetDataStore("Contactizer")
Contactizer.Connect()
```

## Creating and updating objects

`ghost` was designed and simplicity in mind, therefore, there is only one function to save or update an object and that is `ghost.Put()`.

`ghost.Put()` received an object, a pointer to an object, a slice of objects or a pointer to slice of objects.

If the `Id` of an object is `0` then `ghost.Put()` saves a new object, otherwise it updates the object based on it's `Id` value.

If the object passed to `ghost.Put()` is a pointer to an object or a pointer to slice of objects, then the new object get updated with the `Id` value from the database.

If a pointer of a slice of objects is passed to `ghost.Put()`, then `ghost` will update the objects that have `Id`, and save the ones which don't, updating them with their `Id` values from the database.

In code, to save a new object and gets its `Id`

```
Contactizer, _ := gohst.GetDataStore("Contactizer")
contact := Content{}
contact.FirstName = "Ali"
contact.LastName = "Anwar"
err = Contactizer.Put(&contact)

// If it's not necessary to get an object Id when saving a new one
err = Contactizer.Put(contact)
```

The same code would be used to update objects

```
Contactizer, _ := gohst.GetDataStore("Contactizer")

// Get some objects
var contacts []Contact
err = Contactizer.Get(&contacts, []int64{1, 2, 3})

// Change them
contact[0].Country = "Switzerland"
contact[1].Country = "USA"
contact[2].Country = "Germany"

err = Contactizer.Put(contacts)
``` 

**As a warning, it's very important to realise that the object has to be pulled first from the database to be updated, if you assign an object the **`Id`** of another one, you basically overwriting the other one!**

```
// This will empty contact with Id=1 of it's content!!!
contact := Contact{}
contact.Id = 1
err = Contactizer.Put(contact)
```

Passing a slice can create or update multiple objects, depending if they have their `Id` set or not.

##Retrieve objects

`ghost` has three get function, these are `gohst.Get()`, `gohst.GetAll()`, `gohst.GetRaw()`, these are not the only way to retrieve objects from the data store, But the standard `gohst` ones.

### GetAll()

The simples one is `gohst.GetAll()`, it just retrieves all the objects in a table, a very expensive task!

```
Contactizer, _ := gohst.GetDataStore("Contactizer")
var contacts []Contact
err = Contactizer.GetAll(&contacts)
```
It runs a `select * from json_contacts;` and build all the objects from json to Go objects, so use it if you need really really it!

### Get()

The signature for `gohst.Get()` is

```
Get(object interface{}, request interface{})
```

where `object` is always a pointer to a slice of objects where the result will be returned, and `request` is either a slice of type `[]int64` which contains the ids of the object to be retrieved, or an object of type `ghost.Request`.

```
var contacts []Contact
err = Contactizer.Get(&contacts, []int64{1, 2, 3})
```

if the `Id` refers to a deleted object, then no error will be returned, this function is a call to `select * from json_contacts where id in (1,2,3);`

Now, it's not always the case when we have the `id` of objects, maybe we want to search for objects with a certain country or tag in it's tags list, for these an object of type `gohst.Request` is passed.

This example does the same thing as the previous one, it retrieves a list of objects by their `Id` but this time using a `ghost.Request`

```
var contacts []Contact
ids := []int64{1,2,3}
request := &gohst.RequestChain{}
request.Where(gohst.Clause{"Id", "IN", ids})
err := Contactizer.Get(&contacts, request)
```

But really the `ghost.Request` is for queries like this one

```
date, err := time.Parse("2006-01-02 15:04:05", "2011-05-01 00:00:00")
condition := gohst.Clause{"archived_at", "<", date}
request.Where(condition)
var contacts []Contact
Contactizer.Get(&contacts, request)

// Or
request.And(gohst.Clause{"categories", "@>", "$1"})
Contactizer.Get(&contacts, request)
```
**One Last important note about **`Get()`** function, it appends the passed slice, this means when the passed slice is not empty, it will add to it instead of replacing it.**

### GetRaw()

`gohst.GetRaw()` get called the same way as `gohst.Get()` put returns a json string instead, so it's cheaper than `gohst.Get()` since it doesn't unpack the json into a slice of objects. It's useful when only the json list is needed.

## Delete Objects

`gohst.Delete()` works just like `ghost.Get()` by either a slice of `id` or a `ghost.Request`, but it also works on a slice of objects.

```



## Transactions


Transactions are started by calling gohst.Begin() on a datastore. This will not only return a transaction object, but also registers the transaction in the datastore with a name and start time, which allow to call open transactions from several function before committing them, maybe this is a bad idea, just don't use it if you can't live with the consequences.

```
// Begin a transaction
trx := Contactizer.Begin("Change first 3 contacts")
ids := []int64{1,2,3}
var contacts []Contact
Contactizer.Get__(&contacts,ids,trx)
for _, contact := range contacts {
	contact.IsCompany = true
}
Contactizer.Put__(contacts,trx)
Contactizer.Commit(trx)
```


