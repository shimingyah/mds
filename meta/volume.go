package meta

import (
	"context"

	"github.com/shimingyah/mds/pb"
)

// CreateVolume create a volume instance
func (m *Meta) CreateVolume() error {
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
