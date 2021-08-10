package v3

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var db *badger.DB

func Init(path string, inMemory bool) (err error) {
	var opt badger.Options

	if inMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(path)
	}

	//Badger obtains a lock on the directories so multiple processes cannot open the same database at the same time
	if badgerDb, err := badger.Open(opt); err != nil {
		return err
	} else {
		db = badgerDb
	}

	return
}

func Close() error {
	if db != nil {
		return db.Close()
	} else {
		return errors.New("db is not initialized")
	}
}
