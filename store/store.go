package store

import "github.com/shimingyah/mds/pb"

// MergeFunc accepts two byte slices, one representing an existing value, and
// another representing a new value that needs to be ‘merged’ into it. MergeFunc
// contains the logic to perform the ‘merge’ and return an updated value.
// MergeFunc could perform operations like integer addition, list appends etc.
// Note that the ordering of the operands is maintained.
type MergeFunc func(existingVal, newVal []byte) []byte

// Txn for store
type Txn interface {
	Get(key interface{}) (interface{}, error)
	Gets(keys ...interface{}) ([]interface{}, error)
	Set(key, value interface{}) error
	Delete(key interface{}) error
	Deletes(keys ...[]byte) error
	Update(key, value interface{}, columns ...string) error
	Incr(key interface{}, value uint64) (uint64, error)
	Append(key, value interface{}) (interface{}, error)
	Exist(key interface{}) (bool, error)
	Add(key, value interface{}) error
	Merge(key interface{}, fn MergeFunc) (interface{}, error)
	ScanRange(begin, end interface{}) (map[string]interface{}, error)
	ScanKeys(prefix interface{}, limit int) ([]interface{}, error)
	ScanValues(prefix interface{}, limit int) (map[string]interface{}, error)
}

// Formatter serialize/unserialize key/value
type Formatter interface {
	// keys
	CounterKey(key string) interface{}
	InodeKey(volumeID uint32, nodeID uint64) interface{}
	DentryKey(volumeID uint32, parent uint64, name string)
	ChunkKey(volumeID uint32, nodeID uint64, index uint32)
	SymKey(volumeID uint32, nodeID uint64)
	XattrKey(volumeID uint32, nodeID uint64, name string)
	DelFileKey(volumeID uint32, nodeID, length uint64)
	FlockKey(volumeID uint32, nodeID uint64)
	PlockKey(volumeID uint32, nodeID uint64)
	SessionHeartbeatKey(volumeID uint32, sessionID uint64)
	SessionInfoKey(volumeID uint32, sessionID uint64)

	// values
	MarshalAttr(attr *pb.Attr) interface{}
	UnmarshalAttr(buf interface{}, attr *pb.Attr)
	MarshalXAttr()
	UnMarshalXAttr(buf interface{}) ([]byte, []byte)
	MarshalDentry(typ uint8, nodeID uint64) interface{}
	UnmarshalDentry(buf interface{}) (uint8, uint64)
	MarshalSlice(slice *pb.Slice) []byte
	UnmarshalSlice(buf []byte) *pb.Slice
}

// Engine support txn
type Engine interface {
	// Name engine name
	Name() string

	// Txn read-write transaction
	Txn(fn func(Txn) error) error

	// View read-only transaction
	View(fn func(Txn) error) error

	// Update read-write transaction
	Update(fn func(Txn) error) error

	// Snapshot return a mvcc snapshot
	Snapshot() (Snapshot, error)

	// Formatter return a key/value formatter
	Formatter() Formatter

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

// NewEngine return a kv/sql engine instance
func NewEngine() (Engine, error) {
	return newKVEngine("")
}

// IsKeyNotFound return true is there is no such key
func IsKeyNotFound(err error) bool {
	return isKVKeyNotFound(err)
}
