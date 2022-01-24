package meta

import (
	"context"
	"syscall"
)

func (m *meta) Flock(ctx context.Context, volumeID uint32, nodeID uint64, owner uint64, ltype uint32, block bool) syscall.Errno {
	return 0
}

func (m *meta) Getlk(ctx context.Context, volumeID uint32, nodeID uint64, owner uint64, ltype *uint32, start, end *uint64, pid *uint32) syscall.Errno {
	return 0
}

func (m *meta) Setlk(ctx context.Context, volumeID uint32, nodeID uint64, owner uint64, block bool, ltype uint32, start, end uint64, pid uint32) syscall.Errno {
	return 0
}
