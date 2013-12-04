gohst
=====

DataStore Interface


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


