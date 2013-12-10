gohst
=====

DataStore Interface
====================

gohst started as an ORM project, soon it seemed that my project don't need an ORM, they just needed a dependable object storage, something that can have open schema, and support transaction, looking into the NoSQL world I couldn't find something that fits the way I want things to work, so gohst turned from ambitious attempt to a practical tool.

Currently gohst works only with PostgreSQL, the intension was to support different databases, but since a decision was made to use PostgreSQL json extension, most of the code currently assumes that it's PostgreSQL that gohst is talking to, in the future this may change, and the interface may support this.

gohst also has been benchmarked on different laptops, on my MacBook Air, I'm able to pull avg. 3500 object/sec, medium sized, it depends on reflect and encoding/json packages, now with the release of GO 1.2 the json part should be faster.

Objects and Data types
======================

gohst is really simple, it convert a Go object to json, warps it with a predefined record and store it in a table with the same name as the object type. it only expect that the object type defines 3 field

```
type Contact struct {
	Id             int64     `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
```

The `json:"-"` is not necessary, actually if you intend to use Raw queries and forward the json results directly to another server or browser, you may want to remove them, this will include these field inside the json string.

Now, 


User of Transactions
====================

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


