package v3

import (
	badger "github.com/dgraph-io/badger/v3"
	"time"
)

type transaction struct {
	db *badger.DB
	txn *badger.Txn
}

// Set set the key value pair, ttl is optional and is based on millisecond
func (t *transaction) Set(key, value []byte, ttl int64) error {
	e := badger.NewEntry(key, value)
	if ttl > 0 {
		e.WithTTL(time.Duration(ttl) * time.Millisecond)
	}

	return t.txn.SetEntry(e)
}

// Get returns a copy of the value of the item from the value log, writing it to dst slice.
// If nil is passed, or capacity of dst isn't sufficient, a new slice would be allocated and
// returned. Tip: It might make sense to reuse the returned slice as dst argument for the next call.
//
// This function is useful in long running iterate/update transactions to avoid a write deadlock.
// See Github issue: https://github.com/dgraph-io/badger/issues/315
func (t *transaction) Get(key[]byte) (value []byte, err error) {
	item, err := t.txn.Get(key)
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

// ExpiresAt returns a Unix time value indicating when the item will be
// considered expired. 0 indicates that the item will never expire.
func (t *transaction) ExpiresAt(key[]byte) (expiresAt uint64, err error) {
	item, err := t.txn.Get(key)
	if err != nil {
		return
	}

	return item.ExpiresAt(), nil
}

// Delete deletes a key.
//
// This is done by adding a delete marker for the key at commit timestamp.  Any
// reads happening before this timestamp would be unaffected. Any reads after
// this commit would see the deletion.
//
// The current transaction keeps a reference to the key byte slice argument.
// Users must not modify the key until the end of the transaction.
func (t *transaction) Delete(key[]byte) error {
	err := t.txn.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

// Commit commits the transaction, following these steps:
//
// 1. If there are no writes, return immediately.
//
// 2. Check if read rows were updated since txn started. If so, return ErrConflict.
//
// 3. If no conflict, generate a commit timestamp and update written rows' commit ts.
//
// 4. Batch up all writes, write them to value log and LSM tree.
//
// 5. If callback is provided, Badger will return immediately after checking
// for conflicts. Writes to the database will happen in the background.  If
// there is a conflict, an error will be returned and the callback will not
// run. If there are no conflicts, the callback will be called in the
// background upon successful completion of writes or any error during write.
//
// If error is nil, the transaction is successfully committed. In case of a non-nil error, the LSM
// tree won't be updated, so there's no need for any rollback.
func (t *transaction) Commit() error {
	return t.txn.Commit()
}


// Discard discards a created transaction. This method is very important and must be called. Commit
// method calls this internally, however, calling this multiple times doesn't cause any issues. So,
// this can safely be called via a defer right when transaction is created.
//
// NOTE: If any operations are run on a discarded transaction, ErrDiscardedTxn is returned.
func (t *transaction) Discard() {
	t.txn.Discard()
}