package meta

import (
	"context"

	"github.com/shimingyah/mds/pb"
)

func (m *meta) ListSlices(ctx context.Context, volumeID uint32, lastAppend int64, limit int, slices *[]pb.Slice) error {
	return nil
}

func (m *meta) FindCompactChunk(ctx context.Context, volumeID uint32, limit int, lastTime int64) (*[]pb.Chunk, error) {
	return nil, nil
}

func (m *meta) CompactSlices(ctx context.Context, chunk *pb.Chunk) error {
	return nil
}

func (m *meta) CommitCompact(ctx context.Context, volumeID uint32, chunk *pb.Chunk) error {
	return nil
}

func (m *meta) UpdateNeedCompact(chunk *pb.Chunk) error {
	return nil
}

func (m *meta) GetDelFileForVid(ctx context.Context, volumeID uint32, before int64, limit int) error {
	return nil
}

func (m *meta) GetChunks(ctx context.Context, volumeID uint32, nodeID uint64) (chunks []*pb.Chunk, err error) {
	return nil, nil
}

func (m *meta) DelChunk(ctx context.Context, chunk *pb.Chunk) error {
	return nil
}

func (m *meta) DelDelFile(ctx context.Context) error {
	return nil
}
