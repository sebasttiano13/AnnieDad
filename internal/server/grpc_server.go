package server

import (
	"context"
	"net"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano13/AnnieDad/internal/handlers"
	pb "github.com/sebasttiano13/AnnieDad/internal/proto"
	"github.com/sebasttiano13/AnnieDad/internal/service"
	"github.com/sebasttiano13/AnnieDad/pkg/clients"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"google.golang.org/grpc"
)

// GRPSServerSettings stores server settings
type GRPSServerSettings struct {
	SecretKey     string
	CertFile      string
	CertKey       string
	TokenDuration time.Duration
}

// GRPSServer implements gRPC server.
type GRPSServer struct {
	srv *grpc.Server
}

// NewGRPSServer init all interceptors, load TLS, registers services and returns GRPSServer
func NewGRPSServer(
	settings *GRPSServerSettings,
	auth service.AuthRepo,
	media service.MediaRepo,
	s3 *clients.S3Client,
) *GRPSServer {

	tlsCredentials, err := loadTLSCredentials(settings.CertFile, settings.CertKey)
	if err != nil {
		logger.Errorf("cannot load TLS credentials: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(handlers.InterceptorLogger(logger.GetDefault())),
		),
	)

	pb.RegisterAuthServer(s, &handlers.AuthServer{
		Auth: service.NewAuthService(auth),
	})

	pb.RegisterMediaServer(s, &handlers.MediaServer{
		Media: service.NewMediaService(media, s3),
	})
	return &GRPSServer{
		srv: s,
	}
}

// Start запускает grpc сервер.
func (s *GRPSServer) Start(addr string) {
	logger.Infof("Running gRPC server on address %s", addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorf("failed to allocate tcp socket for gRPC server: %v", err)
		os.Exit(1)
	}
	if err := s.srv.Serve(listen); err != nil {
		logger.Errorf("failed to start gRPC server: %v", err)
	}
}

// HandleShutdown закрывает grpc сервер.
func (s *GRPSServer) HandleShutdown(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	<-ctx.Done()
	logger.Infof("shutdown signal caught. shutting down gRPC server")

	s.srv.GracefulStop()
	logger.Infof("gRPC server gracefully shutdown")
}
