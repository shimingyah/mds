package server

import (
	"net"
	"net/http"
	"time"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/shimingyah/mds/logrus"
	"github.com/shimingyah/mds/pb"
	"github.com/shimingyah/mds/utils"
)

var logger = logrus.GetLogger("mds")

// MDS metadata server
type MDS struct {
	endpoint     string
	readTimeout  int
	writeTimeout int
}

// NewMDS return a mds instance
func NewMDS(endpoint string, readTimeout, writeTimeout int) *MDS {
	return &MDS{
		endpoint:     endpoint,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

// RunOnlyHTTP only http
func (m *MDS) RunOnlyHTTP() {
	if err := http.ListenAndServe(m.endpoint, m.RegistHTTPRouter()); err != nil {
		logger.Fatalf("mds failed to serve: %v", err)
	}
}

// Run start pixar http and grpc server
func (m *MDS) Run() {
	listener, err := utils.NewTimeoutListener(m.endpoint,
		time.Duration(m.readTimeout)*time.Second,
		time.Duration(m.writeTimeout)*time.Second)
	if err != nil {
		logger.Fatalf("mds failed to init http listener: %v", err)
	}

	c := cmux.New(listener)
	httpl := c.Match(cmux.HTTP1Fast())
	grpcl := c.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	go m.serveHTTP(httpl)
	go m.serveGRPC(grpcl)

	if err := c.Serve(); err != nil {
		logger.Fatalf("mds failed to serve: %v", err)
	}
}

func (m *MDS) serveHTTP(listener net.Listener) {
	router := m.RegistHTTPRouter()
	srv := &http.Server{Handler: router}
	if err := srv.Serve(listener); err != nil {
		logger.Fatalf("failed to serve http: %v", err)
	}
}

func (m *MDS) serveGRPC(listener net.Listener) {
	srv := grpc.NewServer(
		grpc.InitialWindowSize(InitialWindowSize),
		grpc.InitialConnWindowSize(InitialConnWindowSize),
		grpc.MaxSendMsgSize(MaxSendMsgSize),
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    KeepAliveTime,
			Timeout: KeepAliveTimeout,
		}),
	)
	pb.RegisterMDSServer(srv, m)
	reflection.Register(srv)
	if err := srv.Serve(listener); err != nil {
		logger.Fatalf("failed to serve grpc: %v", err)
	}
}
