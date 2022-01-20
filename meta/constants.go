package meta

const (
	RetryNumber      = 3
	SliceTypeSize    = 24
	plockRecordSize  = 24
	MetricMeta       = "rpc.meta"
	MetricMetaCallee = "meta"

	// InodeIDPrefetchBatch from db
	InodeIDPrefetchBatch = 100
	// SliceIDPrefetchBatch from db
	SliceIDPrefetchBatch = 1000
)

const (
	// CompactChunkMsg is a message to compact a chunk in object store.
	CompactChunkMsg = 1000
)

const (
	TypeFile      = 1 // type for regular file
	TypeDirectory = 2 // type for directory
	TypeSymlink   = 3 // type for symlink
	TypeFIFO      = 4 // type for FIFO node
	TypeBlockDev  = 5 // type for block device
	TypeCharDev   = 6 // type for character device
	TypeSocket    = 7 // type for socket
)
