package meta

// Attr struct in mem
type Attr struct {
	Flags     uint8  // Flags
	Type      uint8  //	Type | Mode 为fuse.Attr.Mode(uint32) :0xf000
	Mode      uint16 // 3位SUID位，第10位为SGID位，第9位为sticky位 + 9位rwxrwxrwx
	Uid       uint32 // owner id
	Gid       uint32 // group id
	Atime     int64  // last access time
	Mtime     int64  // last modified time
	Ctime     int64  // last change time for meta
	Atimensec uint32 // nanosecond part of atime
	Mtimensec uint32 // nanosecond part of mtime
	Ctimensec uint32 // nanosecond part of ctime
	Nlink     uint32 // number of links (sub-directories or hardlinks)
	Length    uint64 // length of regular file
	Rdev      uint32 // device number
	Parent    uint64 // inode of parent, only for Directory
	Full      bool   // the attrib utes are completed or not
}

// Entry struct in mem
type Entry struct {
	NodeId uint64
	Name   string
	Attr   *Attr
}

type Chunk struct {
	Vid        int
	NodeId     uint64
	Indx       uint32
	Slices     []byte
	LastAppend int64
}
