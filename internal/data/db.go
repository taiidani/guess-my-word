package data

import (
	"log"

	"github.com/dgraph-io/badger/v2"
)

type (
	// Backend represents a data store for the application, persisting data to disk
	Backend interface {
		Set(string, []byte) error
		Get(string) ([]byte, error)
	}

	// BadgerBackend is a Backend that uses the Badger K/V library to persist data
	BadgerBackend struct {
		Backend
		db *badger.DB
	}
)

// badgerDir is the location on disk where Badger persists its data
const badgerDir = "./badger"

var db Backend

// SetBackend will assign the given backend to the package scope
// It needs to be called at application initialization in order to have a datastore
// ready for I/O calls
func SetBackend(backend Backend) {
	db = backend
}

// NewBadgerBackend returns a new instance of the Badger DB backend
func NewBadgerBackend(opts *badger.Options) (back Backend, teardown func() error) {
	var err error
	ret := BadgerBackend{}

	if opts == nil {
		defaults := badger.DefaultOptions(badgerDir)
		opts = &defaults
	}

	if ret.db, err = badger.Open(*opts); err != nil {
		log.Panicf("Unable to open Badger database %s", err)
	}

	return &ret, ret.db.Close
}

// Get will retrieve a given key/value pair from the backend
func (b *BadgerBackend) Get(key string) ([]byte, error) {
	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		data, err = item.ValueCopy(nil)
		return err
	})

	return data, err
}

// Set will set a given key/value pair into the backend
func (b *BadgerBackend) Set(key string, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)
		return err
	})
}
