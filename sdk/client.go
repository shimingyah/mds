package sdk

import (
	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/pool"
)

// MDSClient client structure
type MDSClient struct {
	endpoint string
	p        pool.Pool
}

// NewMDSClient return mds instance
func NewMDSClient(endpoint string) (*MDSClient, error) {
	opt := pool.DefaultOptions
	opt.MaxIdle = 2
	opt.MaxActive = 6
	p, err := pool.New(endpoint, opt)
	if err != nil {
		return nil, err
	}

	return &MDSClient{
		endpoint: endpoint,
		p:        p,
	}, nil
}

// Close the mds client
func (c *MDSClient) Close() error {
	return c.p.Close()
}

// PbClient used communication with mds server
// newPbClient and close must be paired.
type PbClient struct {
	pb.MDSClient
	conn pool.Conn
}

// NewPbClient return a pb client instance
func (c *MDSClient) NewPbClient() (*PbClient, error) {
	conn, err := c.p.Get()
	if err != nil {
		return nil, err
	}

	client := pb.NewMDSClient(conn.Value())

	return &PbClient{
		MDSClient: client,
		conn:      conn,
	}, nil
}

// Close release grpc conn
func (c *PbClient) Close() error {
	return c.conn.Close()
}
