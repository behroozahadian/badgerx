package v3

type Badgerx interface {
	Set(key, value []byte, ttl int64) error
	Get(key []byte) (value []byte, err error)
	Delete(key []byte) error
	ExpiresAt(key []byte) (expiresAt uint64, err error)
	Close() error
	NewTransaction(update bool) Txn
}

type Txn interface {
	Set(key, val []byte, ttl int64) error
	Get(key []byte) (value []byte, err error)
	ExpiresAt(key []byte) (expiresAt uint64, err error)
	Delete(key []byte) error

	Discard()
	Commit() error
}
