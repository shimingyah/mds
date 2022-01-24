package store

import (
	"bytes"
	"encoding/binary"

	"github.com/dgraph-io/badger/v3"
)

// MergeFunc accepts two byte slices, one representing an existing value, and
// another representing a new value that needs to be ‘merged’ into it. MergeFunc
// contains the logic to perform the ‘merge’ and return an updated value.
// MergeFunc could perform operations like integer addition, list appends etc.
// Note that the ordering of the operands is maintained.
type MergeFunc func(existingVal, newVal []byte) []byte

// KVTxn for store
type KVTxn interface {
	Get(key []byte) ([]byte, error)
	Gets(keys ...[]byte) ([][]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error
	Deletes(keys ...[]byte) error
	Incr(key []byte, value uint64) (uint64, error)
	Append(key, value []byte) ([]byte, error)
	Exist(key []byte) (bool, error)
	Add(key, value []byte) error
	Merge(key []byte, fn MergeFunc) ([]byte, error)
	ScanRange(begin, end []byte) (map[string][]byte, error)
	ScanKeys(prefix []byte, limit int) ([][]byte, error)
	ScanValues(prefix []byte, limit int) (map[string][]byte, error)
}

// Engine support txn
type Engine interface {
	// Name engine name
	Name() string

	// Txn read-write transaction
	Txn(fn func(KVTxn) error) error

	// View read-only transaction
	View(fn func(KVTxn) error) error

	// Update read-write transaction
	Update(fn func(KVTxn) error) error

	// Snapshot return a mvcc snapshot
	Snapshot() (Snapshot, error)

	// Close engine
	Close() error
}

// Snapshot mvcc snapshot
type Snapshot interface {
	Name() string
	Get(key []byte) ([]byte, error)
	ScanRange(begin, end []byte) (map[string][]byte, error)
	ScanKeys(prefix []byte, limit int) ([][]byte, error)
	Close() error
}

type kvTxn struct {
	*badger.Txn
	db *badger.DB
}

type engine struct {
	db *badger.DB
}

type snapshot struct {
	txn    KVTxn
	rawTxn *badger.Txn
}

func newSnapshot(txn *badger.Txn, db *badger.DB) *snapshot {
	return &snapshot{
		txn: &kvTxn{
			Txn: txn,
			db:  db,
		},
		rawTxn: txn,
	}
}

func (s *snapshot) Name() string {
	return "badgerdb-snapshot"
}

func (s *snapshot) Get(key []byte) ([]byte, error) {
	return s.txn.Get(key)
}

func (s *snapshot) ScanRange(begin, end []byte) (map[string][]byte, error) {
	return s.txn.ScanRange(begin, end)
}

func (s *snapshot) ScanKeys(prefix []byte, limit int) ([][]byte, error) {
	return s.txn.ScanKeys(prefix, limit)
}

func (s *snapshot) Close() error {
	s.rawTxn.Discard()
	return nil
}

// NewEngine creates a new txn engine
func NewEngine(dir string) (Engine, error) {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, err
	}
	return &engine{
		db: db,
	}, nil
}

func (e *engine) Name() string {
	return "badgerdb"
}

func (e *engine) Txn(fn func(KVTxn) error) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *engine) View(fn func(KVTxn) error) error {
	return e.db.View(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *engine) Update(fn func(KVTxn) error) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *engine) Snapshot() (Snapshot, error) {
	if e.db.IsClosed() {
		return nil, badger.ErrDBClosed
	}
	txn := e.db.NewTransaction(false)
	return newSnapshot(txn, e.db), nil
}

func (e *engine) Close() error {
	return e.db.Close()
}

func (txn *kvTxn) Get(key []byte) ([]byte, error) {
	item, err := txn.Txn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (txn *kvTxn) Gets(keys ...[]byte) ([][]byte, error) {
	values := make([][]byte, 0, len(keys))
	for _, key := range keys {
		val, err := txn.Get(key)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func (txn *kvTxn) Set(key, value []byte) error {
	return txn.Txn.Set(key, value)
}

func (txn *kvTxn) Delete(key []byte) error {
	return txn.Txn.Delete(key)
}

func (txn *kvTxn) Deletes(keys ...[]byte) error {
	for _, key := range keys {
		err := txn.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (txn *kvTxn) Incr(key []byte, value uint64) (uint64, error) {
	next := uint64(0)
	item, err := txn.Txn.Get(key)

	switch {
	case err == badger.ErrKeyNotFound:
		next = 0
	case err != nil:
		return 0, err
	default:
		err := item.Value(func(v []byte) error {
			next = binary.BigEndian.Uint64(v)
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	next += value
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], next)
	if err := txn.Txn.Set(key, buf[:]); err != nil {
		return 0, err
	}

	return next, nil
}

func (txn *kvTxn) Append(key, value []byte) ([]byte, error) {
	val, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	new := append(val, value...)
	if err := txn.Set(key, new); err != nil {
		return nil, err
	}
	return new, nil
}

func (txn *kvTxn) Exist(key []byte) (bool, error) {
	_, err := txn.Txn.Get(key)
	switch {
	case err == badger.ErrKeyNotFound:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func (txn *kvTxn) Add(key, value []byte) error {
	return txn.Txn.Add(key, value)
}

func (txn *kvTxn) Merge(key []byte, fn MergeFunc) ([]byte, error) {
	return txn.Merge(key, fn)
}

func (txn *kvTxn) ScanRange(begin, end []byte) (map[string][]byte, error) {
	kvs := make(map[string][]byte)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := txn.NewIterator(opts)
	defer it.Close()
	for it.Seek(begin); it.Valid(); it.Next() {
		key := it.Item().Key()
		if end != nil && bytes.Compare(key, end) > 0 {
			return kvs, nil
		}
		val, err := it.Item().ValueCopy(nil)
		if err != nil {
			return kvs, err
		}
		kvs[string(key)] = val
	}
	return nil, nil
}

func (txn *kvTxn) ScanKeys(prefix []byte, limit int) ([][]byte, error) {
	counter := 0
	keys := make([][]byte, 0, 1024)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		keys = append(keys, it.Item().Key())
		counter++
		if limit > 0 && counter >= limit {
			break
		}
	}
	return keys, nil
}

func (txn *kvTxn) ScanValues(prefix []byte, limit int) (map[string][]byte, error) {
	counter := 0
	kvs := make(map[string][]byte)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := txn.NewIterator(opts)
	defer it.Close()
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		key := string(it.Item().Key())
		val, err := it.Item().ValueCopy(nil)
		if err != nil {
			return kvs, err
		}
		kvs[key] = val
		counter++
		if limit > 0 && counter >= limit {
			break
		}
	}
	return kvs, nil
}

func IsKeyNotFound(err error) bool {
	return err == badger.ErrKeyNotFound
}
