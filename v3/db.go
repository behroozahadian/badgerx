package v3

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	"time"
)

type badgerx struct {
	db *badger.DB
}

// Open returns a new Badgerx interface
func Open(path string, inMemory bool) (badgerDb Badgerx, err error) {
	b := new(badgerx)

	var opt badger.Options

	if inMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(path)
	}

	//Badger obtains a lock on the directories so multiple processes cannot open the same database at the same time
	if db, err := badger.Open(opt); err != nil {
		return
	} else {
		b.db = db
	}

	return b, nil
}

// Close closes a DB. It's crucial to call it to ensure all the pending updates make their way to
// disk. Calling DB.Close() multiple times would still only close the DB once.
func (b *badgerx) Close() error {
	if b.db != nil {
		return b.db.Close()
	} else {
		return errors.New("db is not initialized")
	}
}

// NewTransaction creates a new transaction. Badger supports concurrent execution of transactions,
// providing serializable snapshot isolation, avoiding write skews. Badger achieves this by tracking
// the keys read and at Commit time, ensuring that these read keys weren't concurrently modified by
// another transaction.
//
// For read-only transactions, set update to false. In this mode, we don't track the rows read for
// any changes. Thus, any long running iterations done in this mode wouldn't pay this overhead.
//
// Running transactions concurrently is OK. However, a transaction itself isn't thread safe, and
// should only be run serially. It doesn't matter if a transaction is created by one goroutine and
// passed down to other, as long as the Txn APIs are called serially.
//
// When you create a new transaction, it is absolutely essential to call
// Discard(). This should be done irrespective of what the update param is set
// to. Commit API internally runs Discard, but running it twice wouldn't cause
// any issues.
//
//  txn := db.NewTransaction(false)
//  defer txn.Discard()
//  // Call various APIs.
func (b *badgerx) NewTransaction(update bool) Txn {
	return &transaction{
		db:  b.db,
		txn: b.db.NewTransaction(update),
	}
}

// Set set the key value pair, ttl is optional and is based on millisecond
func (b *badgerx) Set(key, value []byte, ttl int64) error {
	if b.db == nil {
		return errors.New("db not initialized")
	}

	if err := b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		if ttl > 0 {
			e.WithTTL(time.Duration(ttl) * time.Millisecond)
		}
		err := txn.SetEntry(e)
		return err
	}); err != nil {
		return err
	}

	return nil
}

// Get gets the value of provided key
func (b *badgerx) Get(key []byte) (value []byte, err error) {
	if b.db == nil {
		return nil, errors.New("db not initialized")
	}

	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

// ExpiresAt returns expiration	time of provided key
func (b *badgerx) ExpiresAt(key []byte) (expiresAt uint64, err error) {
	if b.db == nil {
		err = errors.New("db not initialized")
		return
	}

	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return
	}

	return item.ExpiresAt(), nil
}

// Delete removes key/value pair
func (b *badgerx) Delete(key []byte) error {
	if b.db == nil {
		return errors.New("db not initialized")
	}

	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Delete(key)
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}
