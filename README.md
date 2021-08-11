#badgerx

Badgerx is a simple wrapper over [BadgerDB](https://github.com/dgraph-io/badger). BadgerDB is an embeddable, persistent 
and fast key-value (KV) database written in pure Go. It is the underlying database for Dgraph, a fast, distributed graph 
database. It's meant to be a performant alternative to non-Go-based key-value stores like RocksDB. 
It becomes handy for those android or ios developers who want to use such database, here Badgerx among with [gomobile](https://github.com/golang/mobile)
can help developers to access some BadgerDB key Operators(Set,Get,Delete,Transactions...). 

####Usage
```go
import (
	"fmt"
	badgerx "github.com/behroozahadian/badgerx/v3"
	"log"
)

func main() {
	// create or open existing database
	badger, err := badgerx.Open("./tmp/badgerdb", false)
	if err != nil {
		log.Panic(err.Error())
	}

	defer badger.Close()

    // set data, here transactions are implemented under the hood
	err = badger.Set([]byte("hello"), []byte("world"), 0)
	if err != nil {
		log.Println(err.Error())
		return
	}

    // get data
	v, err := badger.Get([]byte("hello"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("get hello %s\n", string(v))

    // delete
	err = badger.Delete([]byte("hello"))
	if err != nil {
		log.Println(err.Error())
		return
	}
    
	// you can also use transactions to perform multiple operations in a single transaction
	tx := badger.NewTransaction(true)

	err = tx.Set([]byte("hello_tx"), []byte("world_tx"), 0)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err.Error())
		return
	}

	vTx, err := badger.Get([]byte("hello_tx"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("get hello_tx %s\n", string(vTx))
}
```