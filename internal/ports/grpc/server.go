package grpc

import (
	"fmt"
	"net"

	"github.com/maghaze/auth/pkg/crypto"
	"github.com/maghaze/auth/pkg/token"
	pb "github.com/maghaze/contracts/pbs/go/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	logger *zap.Logger
	crypto crypto.Crypto
	token  token.Token

	api *grpc.Server
	pb.UnimplementedAuthServer
}

func New(log *zap.Logger, c crypto.Crypto, t token.Token) *server {
	s := &server{logger: log, crypto: c, token: t}

	s.api = grpc.NewServer()
	pb.RegisterAuthServer(s.api, s)

	return s
}

func (s *server) Serve(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.Panic("Error listening on tcp address", zap.Int("port", port), zap.Error(err))
	}

	s.logger.Info("GRPC server starts listening on", zap.Int("port", port))
	if err := s.api.Serve(listener); err != nil {
		s.logger.Fatal("Error serving gRPC server", zap.Error(err))
	}
}
