package meta

import (
	"context"

	"github.com/shimingyah/mds/pb"
)

// NextSlice return the next slice id
func (m *Meta) NextSlice(vid int) (uint64, error) {
	client, err := m.newPbClient()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	req := &pb.NextSliceRequest{
		VolumeId: uint64(vid),
	}
	resp, err := client.NextSlice(context.Background(), req)
	if err := checkError(err, resp.Error); err != nil {
		return 0, err
	}
	return resp.SliceId, nil
}

// NewSession return a new session
func (m *Meta) NewSession() error {
	client, err := m.newPbClient()
	if err != nil {
		return err
	}
	defer client.Close()

	req := &pb.CreateVolumeRequest{}
	resp, err := client.CreateVolume(context.Background(), req)
	if err := checkError(err, resp.Error); err != nil {
		return err
	}
	return nil
}
