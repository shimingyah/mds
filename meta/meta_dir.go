package meta

import (
	"context"
	"syscall"

	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/mds/store"
)

func (m *meta) StatFS(ctx context.Context, volumeID uint32, totalspace, availspace, iused, iavail *uint64) syscall.Errno {
	*totalspace = 0
	*availspace = 0
	*iused = 0
	*iavail = 0
	return 0
}

func (m *meta) Lookup(ctx context.Context, volumeID uint32, parent uint64, name string, attr *pb.Attr) syscall.Errno {
	val, err := m.Get(m.DentryKey(volumeID, parent, name))
	if store.IsKeyNotFound(err) {
		return syscall.ENOENT
	}
	if err != nil {
		return errno(err)
	}
	_, nodeID := m.unmarshalDentry(val)
	val, err = m.Get(m.InodeKey(volumeID, nodeID))
	if store.IsKeyNotFound(err) {
		return syscall.ENOENT
	}
	if err != nil {
		return errno(err)
	}
	m.unmarshalAttr(val, attr)
	attr.Inode = nodeID
	return 0
}

func (m *meta) Readdir(ctx context.Context, volumeID uint32, nodeID uint64, wantattr uint8, off, size int, seek string, entries *[]*pb.Entry) syscall.Errno {
	if off < 0 || size < 1 {
		return syscall.EINVAL
	}

	*entries = make([]*pb.Entry, 0, size)
	// .  .. get dir and parent dir attr
	if off < 2 {
		dirEntry := make([]*pb.Entry, 0, 2)

		// get dir attr
		var attr pb.Attr // .
		errno := m.GetAttr(ctx, volumeID, nodeID, &attr)
		if errno != 0 {
			return errno
		}

		// get parent
		var pattr pb.Attr // ..
		if nodeID == 1 {
			pattr = attr
		} else {
			errno = m.GetAttr(ctx, volumeID, attr.Parent, &pattr)
			if errno != 0 {
				return errno
			}
		}
		dirEntry = append(dirEntry, &pb.Entry{
			Inode: nodeID,
			Name:  ".",
			Attr:  &attr,
		})
		dirEntry = append(dirEntry, &pb.Entry{
			Inode: attr.Parent,
			Name:  "..",
			Attr:  &pattr,
		})

		for i := off; i < 2 && size > 0; i++ {
			*entries = append(*entries, dirEntry[i])
			off--
			size--
		}

		if off < 0 {
			off = 0
		}
		if size <= 0 {
			return 0
		}
	} else {
		off -= 2
	}

	key := m.DentryKey(volumeID, nodeID, seek)
	vals, err := m.ScanValues(key, size)
	if err != nil {
		return errno(err)
	}
	prefix := len(key)
	for name, buf := range vals {
		typ, inode := m.unmarshalDentry(buf)
		entry := &pb.Entry{
			Inode: inode,
			Name:  string([]byte(name)[prefix:]),
			Attr:  &pb.Attr{Type: uint32(typ)},
		}

		if wantattr != 0 {
			val, err := m.Get(m.InodeKey(volumeID, inode))
			if store.IsKeyNotFound(err) {
				continue
			}
			m.unmarshalAttr(val, entry.Attr)
		}

		*entries = append(*entries)
	}

	return 0
}

func (m *meta) Resolve(ctx context.Context, parent uint64, path string, attr *pb.Attr) syscall.Errno {
	return 0
}
