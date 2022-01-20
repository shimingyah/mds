package meta

import (
	"sync/atomic"

	"github.com/shimingyah/mds/store"
)

type meta struct {
	engine    store.Engine
	newSpace  int64
	newInodes int64
}

func (m *meta) Txn(fn func(store.KVTxn) error) error {
	return m.engine.Txn(fn)
}

func (m *meta) updateStats(space int64, inodes int64) {
	atomic.AddInt64(&m.newSpace, space)
	atomic.AddInt64(&m.newInodes, inodes)
}

// SliceMergeFunc Merge function to append one byte slice to another
func SliceMergeFunc(originalValue, newValue []byte) []byte {
	return append(originalValue, newValue...)
}
