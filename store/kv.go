package store

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type kvEngine struct {
	db *badger.DB
}

type kvTxn struct {
	*badger.Txn
	db *badger.DB
}

type kvSnapshot struct {
	txn    Txn
	rawTxn *badger.Txn
}

func newKVEngine(dir string) (Engine, error) {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, err
	}
	return &kvEngine{
		db: db,
	}, nil
}

func (e *kvEngine) Name() string {
	return "kv-engine"
}

func (e *kvEngine) Txn(fn func(Txn) error) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *kvEngine) View(fn func(Txn) error) error {
	return e.db.View(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *kvEngine) Update(fn func(Txn) error) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return fn(&kvTxn{txn, e.db})
	})
}

func (e *kvEngine) Snapshot() (Snapshot, error) {
	if e.db.IsClosed() {
		return nil, badger.ErrDBClosed
	}
	txn := e.db.NewTransaction(false)
	return newKVSnapshot(txn, e.db), nil
}

func (e *kvEngine) Close() error {
	return e.db.Close()
}

func newKVSnapshot(txn *badger.Txn, db *badger.DB) *kvSnapshot {
	return &kvSnapshot{
		txn: &kvTxn{
			Txn: txn,
			db:  db,
		},
		rawTxn: txn,
	}
}

func (s *kvSnapshot) Name() string {
	return "kv-snapshot"
}

func (s *kvSnapshot) Get(key []byte) ([]byte, error) {
	return s.txn.Get(key)
}

func (s *kvSnapshot) ScanRange(begin, end []byte) (map[string][]byte, error) {
	return s.txn.ScanRange(begin, end)
}

func (s *kvSnapshot) ScanKeys(prefix []byte, limit int) ([][]byte, error) {
	return s.txn.ScanKeys(prefix, limit)
}

func (s *kvSnapshot) Close() error {
	s.rawTxn.Discard()
	return nil
}

func (t *kvTxn) Get(key []byte) ([]byte, error) {
	item, err := t.Txn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (t *kvTxn) Gets(keys ...[]byte) ([][]byte, error) {
	values := make([][]byte, 0, len(keys))
	for _, key := range keys {
		val, err := t.Get(key)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func (t *kvTxn) Set(key, value []byte) error {
	return t.Txn.Set(key, value)
}

func (t *kvTxn) Delete(key []byte) error {
	return t.Txn.Delete(key)
}

func (t *kvTxn) Deletes(keys ...[]byte) error {
	for _, key := range keys {
		err := t.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *kvTxn) Update(key, value []byte, columns ...string) error {
	return t.Txn.Set(key, value)
}

func (t *kvTxn) Incr(key []byte, value uint64) (uint64, error) {
	next := uint64(0)
	item, err := t.Txn.Get(key)

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
	if err := t.Txn.Set(key, buf[:]); err != nil {
		return 0, err
	}

	return next, nil
}

func (t *kvTxn) Append(key, value []byte) ([]byte, error) {
	val, err := t.Get(key)
	if err != nil {
		return nil, err
	}

	new := append(val, value...)
	if err := t.Set(key, new); err != nil {
		return nil, err
	}
	return new, nil
}

func (t *kvTxn) Exist(key []byte) (bool, error) {
	_, err := t.Txn.Get(key)
	switch {
	case err == badger.ErrKeyNotFound:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func (t *kvTxn) Add(key, value []byte) error {
	return t.Txn.SetEntry(badger.NewEntry(key, value))
}

func (t *kvTxn) Merge(key []byte, fn MergeFunc) ([]byte, error) {
	m := t.db.GetMergeOperator(key, badger.MergeFunc(fn), 5*time.Second)
	defer m.Stop()
	return m.Get()
}

func (t *kvTxn) ScanRange(begin, end []byte) (map[string][]byte, error) {
	kvs := make(map[string][]byte)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := t.NewIterator(opts)
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

func (t *kvTxn) ScanKeys(prefix []byte, limit int) ([][]byte, error) {
	counter := 0
	keys := make([][]byte, 0, 1024)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := t.NewIterator(opts)
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

func (t *kvTxn) ScanValues(prefix []byte, limit int) (map[string][]byte, error) {
	counter := 0
	kvs := make(map[string][]byte)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := t.NewIterator(opts)
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

func isKVKeyNotFound(err error) bool {
	return err == badger.ErrKeyNotFound
}
