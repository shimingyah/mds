package meta

import (
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/shimingyah/mds/logrus"
	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/mds/store"
	"github.com/shimingyah/mds/utils"
)

var logger = logrus.GetLogger("mds")

func (m *meta) keyLen(args ...interface{}) (length int) {
	for _, arg := range args {
		switch arg := arg.(type) {
		case byte:
			length++
		case uint32:
			length += 4
		case uint64:
			length += 8
		case string:
			length += len(arg)
		default:
			logger.Fatalf("invalid type: %T, value: %v", arg, arg)
		}
	}
	return length
}

func (m *meta) fmtKey(args ...interface{}) []byte {
	b := utils.NewBuffer(uint32(m.keyLen(args...)))
	for _, arg := range args {
		switch arg := arg.(type) {
		case byte:
			b.Put8(arg)
		case uint32:
			b.Put32(arg)
		case uint64:
			b.Put64(arg)
		case string:
			b.Put([]byte(arg))
		default:
			logger.Fatalf("invalid type: %T, value: %v", arg, arg)
		}
	}
	return b.Bytes()
}

//   ConfKeys   ---> C
//   MetaKeys 	---> M
//	 VolumeID   ---> vvvv
//   InodeID 	---> iiiiiiii
//   Length 	---> llllllll
//   Index 		---> nnnn
//   name 		---> ...
//   chunkID 	---> cccccccc
//   SessionID	---> ssssssss
//
// All Conf keys:
//
//
// All meta keys:
//   MC...               	 counter
//   MIvvvviiiiiiiiA         inode attribute
//   MIvvvviiiiiiiiD...      dentry
//   MIvvvviiiiiiiiCnnnn     file chunks
//   MIvvvviiiiiiiiS         symlink target
//   MIvvvviiiiiiiiX...      extented attribute
//   MDvvvviiiiiiiillllllll  delete inodes
//   MFvvvviiiiiiii          Flocks
//   MPvvvviiiiiiii          POSIX locks
//   MSHvvvvssssssss         session heartbeat
//   MSIvvvvssssssss         session info

func (m *meta) CounterKey(key string) []byte {
	return m.fmtKey("MC", key)
}

func (m *meta) InodeKey(volumeID uint32, nodeID uint64) []byte {
	return m.fmtKey("MI", volumeID, nodeID, "A")
}

func (m *meta) DentryKey(volumeID uint32, parent uint64, name string) []byte {
	return m.fmtKey("MI", volumeID, parent, "D", name)
}

func (m *meta) ChunkKey(volumeID uint32, nodeID uint64, index uint32) []byte {
	return m.fmtKey("MI", volumeID, nodeID, "C", index)
}

func (m *meta) SymKey(volumeID uint32, nodeID uint64) []byte {
	return m.fmtKey("MI", volumeID, nodeID, "S")
}

func (m *meta) XattrKey(volumeID uint32, nodeID uint64, name string) []byte {
	return m.fmtKey("MI", volumeID, nodeID, "X", name)
}

func (m *meta) DelFileKey(volumeID uint32, nodeID, length uint64) []byte {
	return m.fmtKey("MD", volumeID, nodeID, length)
}

func (m *meta) FlockKey(volumeID uint32, nodeID uint64) []byte {
	return m.fmtKey("MF", volumeID, nodeID)
}

func (m *meta) PlockKey(volumeID uint32, nodeID uint64) []byte {
	return m.fmtKey("MP", volumeID, nodeID)
}

func (m *meta) SessionHeartbeatKey(volumeID uint32, sessionID uint64) []byte {
	return m.fmtKey("MSH", volumeID, sessionID)
}

func (m *meta) SessionInfoKey(volumeID uint32, sessionID uint64) []byte {
	return m.fmtKey("MSI", volumeID, sessionID)
}

func (m *meta) Get(key []byte) (val []byte, err error) {
	err = m.engine.View(func(txn store.Txn) error {
		val, err = txn.Get(key)
		return err
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (m *meta) GetCounter(key []byte) (val uint64, err error) {
	err = m.engine.Update(func(txn store.Txn) error {
		val, err = txn.Incr(key, 0)
		return err
	})
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (m *meta) ScanKeys(prefix []byte, limit int) (keys [][]byte, err error) {
	err = m.engine.View(func(txn store.Txn) error {
		keys, err = txn.ScanKeys(prefix, limit)
		return err
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (m *meta) ScanValues(prefix []byte, limit int) (values map[string][]byte, err error) {
	err = m.engine.View(func(txn store.Txn) error {
		values, err = txn.ScanValues(prefix, limit)
		return err
	})
	return values, err
}

func (m *meta) marshalAttr(attr *pb.Attr) []byte {
	w := utils.NewBuffer(36 + 24 + 4 + 16)
	w.Put8(uint8(attr.Flags))
	w.Put16((uint16(attr.Type) << 12) | (uint16(attr.Mode) & 0xFFF))
	w.Put32(attr.Uid)
	w.Put32(attr.Gid)
	w.Put64(attr.Atime)
	w.Put64(attr.Mtime)
	w.Put64(attr.Ctime)
	w.Put32(attr.Atimensec)
	w.Put32(attr.Mtimensec)
	w.Put32(attr.Ctimensec)
	w.Put32(attr.Nlink)
	w.Put64(attr.Length)
	w.Put32(attr.Rdev)
	w.Put64(attr.Inode)
	w.Put64(attr.Parent)
	return w.Bytes()
}

func (m *meta) unmarshalAttr(buf []byte, attr *pb.Attr) {
	if attr == nil {
		return
	}
	b := utils.FromBuffer(buf)
	attr.Flags = uint32(b.Get8())
	attr.Mode = uint32(b.Get16())
	attr.Type = uint32(attr.Mode >> 12)
	attr.Mode &= 0xFFF
	attr.Uid = b.Get32()
	attr.Gid = b.Get32()
	attr.Atime = b.Get64()
	attr.Mtime = b.Get64()
	attr.Ctime = b.Get64()
	attr.Atimensec = b.Get32()
	attr.Mtimensec = b.Get32()
	attr.Ctimensec = b.Get32()
	attr.Nlink = b.Get32()
	attr.Length = b.Get64()
	attr.Rdev = b.Get32()
	attr.Inode = b.Get64()
	attr.Parent = b.Get64()
	attr.Full = true
}

func (m *meta) marshalDentry(typ uint8, nodeID uint64) []byte {
	b := utils.NewBuffer(9)
	b.Put8(typ)
	b.Put64(nodeID)
	return b.Bytes()
}

func (m *meta) unmarshalDentry(buf []byte) (uint8, uint64) {
	b := utils.FromBuffer(buf)
	return b.Get8(), b.Get64()
}

func (m *meta) marshalSlice(slice *pb.Slice) []byte {
	w := utils.NewBuffer(SliceTypeSize)
	w.Put64(slice.Id)
	w.Put32(slice.Size_)
	w.Put32(slice.Pos)
	w.Put32(slice.Off)
	w.Put32(slice.Len)
	return w.Bytes()
}

func (m *meta) unmarshalSlice(buf []byte) *pb.Slice {
	b := utils.ReadBuffer(buf)
	return &pb.Slice{
		Id:    b.Get64(),
		Size_: b.Get32(),
		Pos:   b.Get32(),
		Off:   b.Get32(),
		Len:   b.Get32(),
	}
}

func (m *meta) unmarshalSlices(buf []byte) []*pb.Slice {
	if buf == nil {
		return nil
	}
	var sliceList []*pb.Slice
	for i := 0; i < len(buf); i += SliceTypeSize {
		slice := m.unmarshalSlice(buf[i : i+SliceTypeSize])
		sliceList = append(sliceList, slice)
	}
	return sliceList
}

func (m *meta) shouldCompact(slices []*pb.Slice) bool {
	sliceCompactLimit := 100
	if sliceCompactLimit == 0 {
		return false
	}
	if (len(slices)/SliceTypeSize)%sliceCompactLimit == sliceCompactLimit-1 {
		return true
	}
	return false
}

func errno(err error) syscall.Errno {
	if err == nil {
		return 0
	}
	if eno, ok := err.(syscall.Errno); ok {
		return eno
	}
	if strings.HasPrefix(err.Error(), "OOM") {
		return syscall.ENOSPC
	}
	logger.Errorf("error: %v\n %s", err, debug.Stack())
	return syscall.EIO
}
