package meta

import (
	"errors"

	"github.com/shimingyah/mds/pb"
)

type Meta struct {
	mdsClient *MDSClient
}

func NewMeta() (*Meta, error) {
	return &Meta{}, nil
}

func (m *Meta) newPbClient() (*PbClient, error) {
	return m.mdsClient.NewPbClient()
}

func checkError(err error, pbErr *pb.Error) error {
	if err != nil {
		return err
	}
	if pbErr == nil {
		return nil
	}
	if pbErr.Errcode != 0 {
		return errors.New(pbErr.Errmsg)
	}
	return nil
}
