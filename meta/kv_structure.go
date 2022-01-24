package meta

import (
	"context"
	"syscall"

	"github.com/shimingyah/mds/pb"
)

func (m *meta) Mknod(ctx context.Context, volumeID uint32, parent uint64, name string, typ uint8, mode uint16, umask uint16, rdev uint32, attr *pb.Attr) syscall.Errno {
	return m.mknod(ctx, volumeID, parent, name, typ, mode, umask, rdev, "", attr)
}

func (m *meta) Mkdir(ctx context.Context, volumeID uint32, parent uint64, name string, mode uint16, umask uint16, copysgid uint8, attr *pb.Attr) syscall.Errno {
	return m.Mknod(ctx, volumeID, parent, name, TypeDirectory, mode, umask, 0, attr)
}

func (m *meta) Unlink(ctx context.Context, volumeID uint32, parent uint64, name string) syscall.Errno {
	return 0
}

func (m *meta) Rmdir(ctx context.Context, volumeID uint32, parent uint64, name string) syscall.Errno {
	return 0
}

func (m *meta) Rename(ctx context.Context, volumeID uint32, parentSrc uint64, nameSrc string, parentDst uint64, nameDst string) syscall.Errno {
	return 0
}

func (m *meta) Link(ctx context.Context, volumeID uint32, nodeIdSrc, parent uint64, name string, attr *pb.Attr) syscall.Errno {
	return 0
}

func (m *meta) Symlink(ctx context.Context, volumeID uint32, parent uint64, name string, path string, attr *pb.Attr) syscall.Errno {
	return 0
}

func (m *meta) ReadLink(ctx context.Context, volumeID uint32, nodeID uint64, path *[]byte) syscall.Errno {
	return 0
}

func (m *meta) Access(ctx context.Context, volumeID uint32, nodeID uint64, mode uint8, attr *pb.Attr) syscall.Errno {
	return 0
}

func (m *meta) Forget(volumeID uint32, nodeID, nlookup uint64) {
}

func (m *meta) mknod(ctx context.Context, volumeID uint32, parent uint64, name string, typ uint8, mode uint16, umask uint16, rdev uint32, pointTo string, attr *pb.Attr) syscall.Errno {
	return 0
}
