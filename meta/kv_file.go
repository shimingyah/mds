package meta

import (
	"context"
	"syscall"
	"time"

	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/mds/store"
	"github.com/shimingyah/mds/utils"
)

func (m *meta) Read(ctx context.Context, volumeID uint32, nodeID uint64, indx uint32, slices *[]*pb.Slice) syscall.Errno {
	return errno(m.Txn(func(txn store.KVTxn) error {
		val, err := txn.Merge(m.ChunkKey(volumeID, nodeID, indx), SliceMergeFunc)
		if store.IsKeyNotFound(err) {
			return syscall.ENOENT
		}
		if err != nil {
			return err
		}
		*slices = m.unmarshalSlices(val)
		if m.shouldCompact(*slices) {
			// go compaction
		}
		return nil
	}))
}

func (m *meta) Write(ctx context.Context, volumeID uint32, nodeID uint64, indx uint32, off uint32, slice pb.Slice) syscall.Errno {
	var newSpace int64
	err := m.Txn(func(txn store.KVTxn) error {
		var attr pb.Attr

		val, err := txn.Get(m.InodeKey(volumeID, nodeID))
		if store.IsKeyNotFound(err) {
			return syscall.ENOENT
		}
		m.unmarshalAttr(val, &attr)

		if attr.Type != TypeFile {
			return syscall.EPERM
		}

		newLength := uint64(indx)*uint64(1111) + uint64(off) + uint64(slice.Len)
		if newLength > attr.Length {
			newSpace = utils.Align4K(newLength) - utils.Align4K(attr.Length)
			attr.Length = newLength
		}

		now := time.Now()
		attr.Mtime = uint64(now.Unix())
		attr.Mtimensec = uint32(now.Nanosecond())
		attr.Ctime = uint64(now.Unix())
		attr.Ctimensec = uint32(now.Nanosecond())

		err = txn.Add(m.ChunkKey(volumeID, nodeID, indx), m.marshalSlice(&slice))
		if err != nil {
			return err
		}
		err = txn.Set(m.InodeKey(volumeID, nodeID), m.marshalAttr(&attr))
		if err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		m.updateStats(newSpace, 0)
	}
	return errno(err)
}
