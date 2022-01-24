package meta

import (
	"context"
	"syscall"
	"time"

	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/mds/store"
)

func (m *meta) GetAttr(ctx context.Context, volumeID uint32, nodeID uint64, attr *pb.Attr) syscall.Errno {
	val, err := m.Get(m.InodeKey(volumeID, nodeID))
	if store.IsKeyNotFound(err) {
		return syscall.ENOENT
	}
	if val != nil {
		m.unmarshalAttr(val, attr)
	}
	return errno(err)
}

func (m *meta) SetAttr(ctx context.Context, volumeID uint32, nodeID uint64, uid, gid uint32, modeSet, gidSet, uidSet, atimeSet, mtimeSet, sizeSet bool, attr *pb.Attr) syscall.Errno {
	return errno(m.Txn(func(txn store.KVTxn) error {
		val, err := m.Get(m.InodeKey(volumeID, nodeID))
		if store.IsKeyNotFound(err) {
			return syscall.ENOENT
		}
		if err != nil || val == nil {
			logger.Errorf("[meta] get inode(vid: %d, nodeId: %d) failed, err: %v", volumeID, nodeID, err)
			return err
		}

		cur, changed, now := &pb.Attr{}, false, time.Now()
		m.unmarshalAttr(val, cur)

		if (uidSet || gidSet) && modeSet {
			attr.Mode |= (cur.Mode & 06000)
		}
		if (uidSet || gidSet) && (cur.Mode&06000) != 0 {
			if cur.Mode&01777 != cur.Mode {
				cur.Mode &= 01777
				changed = true
			}
			attr.Mode &= 01777
		}
		if uidSet && cur.Uid != attr.Uid {
			cur.Uid = attr.Uid
			changed = true
		}
		if gidSet && cur.Gid != attr.Gid {
			cur.Gid = attr.Gid
			changed = true
		}
		if modeSet {
			if uid != 0 && (attr.Mode&02000) != 0 {
				if gid != cur.Gid {
					attr.Mode &= 05777
				}
			}
			if attr.Mode != cur.Mode {
				cur.Mode = attr.Mode
				changed = true
			}
		}

		if atimeSet {
			cur.Atime = attr.Atime
			changed = true
		}
		if mtimeSet {
			cur.Mtime = attr.Mtime
			changed = true
		}
		if !changed {
			*attr = *cur
			return nil
		}
		cur.Ctime = uint64(now.UnixNano())
		cur.Ctimensec = uint32(now.Nanosecond())
		err = txn.Set(m.InodeKey(volumeID, nodeID), m.marshalAttr(cur))
		if err != nil {
			logger.Errorf("[meta] update inode(vid: %d, nodeId: %d) failed, err:%v", volumeID, nodeID, err)
			return err
		}
		*attr = *cur
		return nil
	}))
}

func (m *meta) GetXAttr(ctx context.Context, volumeID uint32, nodeID uint64, name string, vbuff *[]byte) syscall.Errno {
	val, err := m.Get(m.XattrKey(volumeID, nodeID, name))
	if store.IsKeyNotFound(err) {
		return ENOATTR
	}
	if val != nil {
		*vbuff = val
	}
	return errno(err)
}

func (m *meta) SetXAttr(ctx context.Context, volumeID uint32, nodeID uint64, name string, value []byte) syscall.Errno {
	if name == "" {
		return syscall.EINVAL
	}
	return errno(m.Txn(func(txn store.KVTxn) error {
		return txn.Set(m.XattrKey(volumeID, nodeID, name), value)
	}))
}

func (m *meta) ListXAttr(ctx context.Context, volumeID uint32, nodeID uint64, names *[]byte) syscall.Errno {
	keys, err := m.ScanKeys(m.XattrKey(volumeID, nodeID, ""), 0)
	if store.IsKeyNotFound(err) {
		return ENOATTR
	}
	if err != nil {
		return errno(err)
	}
	*names = nil
	prefix := len(m.XattrKey(volumeID, nodeID, ""))
	for _, name := range keys {
		*names = append(*names, name[prefix:]...)
		*names = append(*names, 0)
	}
	return 0
}

func (m *meta) RemoveXAttr(ctx context.Context, volumeID uint32, nodeID uint64, name string) syscall.Errno {
	if name == "" {
		return syscall.EINVAL
	}
	return errno(m.Txn(func(txn store.KVTxn) error {
		key := m.XattrKey(volumeID, nodeID, name)
		has, err := txn.Exist(key)
		if err != nil {
			return err
		}
		if !has {
			return ENOATTR
		}
		return txn.Delete(key)
	}))
}
