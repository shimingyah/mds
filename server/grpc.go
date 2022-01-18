package server

import (
	"context"
	"time"

	"github.com/shimingyah/mds/pb"
)

const (
	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 1GB is to provide system's throughput.
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize we set it 1GB is to provide system's throughput.
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 4 << 30
)

// CreateVolume return an volume
func (m *MDS) CreateVolume(context.Context, *pb.CreateVolumeRequest) (*pb.CreateVolumeResponse, error) {
	return nil, nil
}

func (m *MDS) NewSession(context.Context, *pb.NewSessionRequest) (*pb.NewSessionResponse, error) {
	return nil, nil
}

func (m *MDS) SessionHeartbeat(context.Context, *pb.SessionHeartbeatRequest) (*pb.SessionHeartbeatResponse, error) {
	return nil, nil
}

func (m *MDS) CleanInvalidSessions(context.Context, *pb.CleanInvalidSessionsRequest) (*pb.CleanInvalidSessionsResponse, error) {
	return nil, nil
}

func (m *MDS) GetAttr(context.Context, *pb.GetAttrRequest) (*pb.GetAttrResponse, error) {
	return nil, nil
}

func (m *MDS) SetAttr(context.Context, *pb.SetAttrRequest) (*pb.SetAttrResponse, error) {
	return nil, nil
}

func (m *MDS) GetXAttr(context.Context, *pb.GetXAttrRequest) (*pb.GetXAttrResponse, error) {
	return nil, nil
}

func (m *MDS) SetXAttr(context.Context, *pb.SetXAttrResponse) (*pb.SetXAttrResponse, error) {
	return nil, nil
}

func (m *MDS) ListXAttr(context.Context, *pb.ListXAttrRequest) (*pb.ListXAttrResponse, error) {
	return nil, nil
}

func (m *MDS) RemoveXAttr(context.Context, *pb.RemoveXAttrRequest) (*pb.RemoveXAttrResponse, error) {
	return nil, nil
}

func (m *MDS) StatFS(context.Context, *pb.StatFSRequest) (*pb.StatFSResponse, error) {
	return nil, nil
}

func (m *MDS) Mknod(context.Context, *pb.MknodRequest) (*pb.MknodResponse, error) {
	return nil, nil
}

func (m *MDS) Mkdir(context.Context, *pb.MkdirRequest) (*pb.MkdirResponse, error) {
	return nil, nil
}

func (m *MDS) Unlink(context.Context, *pb.UnlinkRequest) (*pb.UnlinkResponse, error) {
	return nil, nil
}

func (m *MDS) Rmdir(context.Context, *pb.RmdirRequest) (*pb.RmdirResponse, error) {
	return nil, nil
}

func (m *MDS) Rename(context.Context, *pb.RenameRequest) (*pb.RenameResponse, error) {
	return nil, nil
}

func (m *MDS) Link(context.Context, *pb.LinkRequest) (*pb.LinkResponse, error) {
	return nil, nil
}

func (m *MDS) Symlink(context.Context, *pb.SymlinkRequest) (*pb.SymlinkResponse, error) {
	return nil, nil
}

func (m *MDS) Access(context.Context, *pb.AccessRequest) (*pb.AccessResponse, error) {
	return nil, nil
}

func (m *MDS) Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error) {
	return nil, nil
}

func (m *MDS) Open(context.Context, *pb.OpenRequest) (*pb.OpenResponse, error) {
	return nil, nil
}

func (m *MDS) Close(context.Context, *pb.CloseRequest) (*pb.CloseResponse, error) {
	return nil, nil
}

func (m *MDS) Read(context.Context, *pb.ReadRequest) (*pb.ReadResponse, error) {
	return nil, nil
}

func (m *MDS) Write(context.Context, *pb.WriteRequest) (*pb.WriteResponse, error) {
	return nil, nil
}

func (m *MDS) Truncate(context.Context, *pb.TruncateRequest) (*pb.TruncateResponse, error) {
	return nil, nil
}

func (m *MDS) Fallocate(context.Context, *pb.FallocateRequest) (*pb.FallocateResponse, error) {
	return nil, nil
}

func (m *MDS) Flock(context.Context, *pb.FlockRequest) (*pb.FlockResponse, error) {
	return nil, nil
}

func (m *MDS) Getlk(context.Context, *pb.GetlkRequest) (*pb.GetlkResponse, error) {
	return nil, nil
}

func (m *MDS) Setlk(context.Context, *pb.SetlkRequest) (*pb.SetlkResponse, error) {
	return nil, nil
}

func (m *MDS) CopyFileRange(context.Context, *pb.CopyFileRangeRequest) (*pb.CopyFileRangeResponse, error) {
	return nil, nil
}

func (m *MDS) Lookup(context.Context, *pb.LookupRequest) (*pb.LookupResponse, error) {
	return nil, nil
}

func (m *MDS) Resolve(context.Context, *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	return nil, nil
}

func (m *MDS) Readdir(context.Context, *pb.ReaddirRequest) (*pb.ReaddirResponse, error) {
	return nil, nil
}

func (m *MDS) NextInode(context.Context, *pb.NextInodeRequest) (*pb.NextInodeResponse, error) {
	return nil, nil
}

func (m *MDS) NextSlice(context.Context, *pb.NextSliceRequest) (*pb.NextSliceResponse, error) {
	return nil, nil
}

func (m *MDS) FindCompactChunk(context.Context, *pb.FindCompactChunkRequest) (*pb.FindCompactChunkResponse, error) {
	return nil, nil
}

func (m *MDS) CompactSlices(context.Context, *pb.CompactSlicesRequest) (*pb.CompactSlicesResponse, error) {
	return nil, nil
}

func (m *MDS) UpdateNeedCompact(context.Context, *pb.UpdateNeedCompactRequest) (*pb.UpdateNeedCompactResponse, error) {
	return nil, nil
}

func (m *MDS) CommitCompact(context.Context, *pb.CommitCompactRequest) (*pb.CommitCompactResponse, error) {
	return nil, nil
}

func (m *MDS) ListSlices(context.Context, *pb.ListSlicesRequest) (*pb.ListSlicesResponse, error) {
	return nil, nil
}

func (m *MDS) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksResponse, error) {
	return nil, nil
}

func (m *MDS) DelChunk(context.Context, *pb.DelChunkRequest) (*pb.DelChunkResponse, error) {
	return nil, nil
}

func (m *MDS) GetDelFiles(context.Context, *pb.GetDelFilesRequest) (*pb.GetDelFilesResponse, error) {
	return nil, nil
}

func (m *MDS) DelDelFile(context.Context, *pb.DelDelFileRequest) (*pb.DelDelFileResponse, error) {
	return nil, nil
}
