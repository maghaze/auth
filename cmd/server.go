package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/maghaze/auth/internal/config"
	"github.com/maghaze/auth/internal/ports/grpc"
	"github.com/maghaze/auth/internal/ports/rest"
	"github.com/maghaze/auth/pkg/crypto"
	"github.com/maghaze/auth/pkg/logger"
	"github.com/maghaze/auth/pkg/token"
)

type Server struct {
	// values from commandline arguments and parameters
	managementPort int
	grpcPort       int

	// values that have been initialized
	config *config.Config
	logger *zap.Logger
	crypto crypto.Crypto
	token  token.Token
}

func (server Server) Command(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Serve the authentication and authorization server",
		Run: func(_ *cobra.Command, _ []string) {
			defer func() {
				<-ctx.Done()
				server.logger.Warn("Got interruption signal, gracefully shutdown the server")
				server.shutdown()
			}()

			server.initialize()
			server.serve()
		},
	}

	cmd.Flags().IntVar(&server.managementPort, "management-port", 8080, "The port the metrics and probe endpoints exposed from")
	cmd.Flags().IntVar(&server.grpcPort, "grpc-port", 9090, "The port the grpc endpoint listens on")

	return cmd
}

func (server *Server) initialize() {
	var err error

	server.config = config.Load(true)
	server.logger = logger.NewZap(server.config.Logger)

	server.crypto = crypto.New(server.config.Crypto)
	server.token, err = token.New(server.config.Token)
	if err != nil {
		server.logger.Panic("Error creating token object", zap.Error(err))
	}
}

func (server *Server) serve() {
	go rest.New(server.logger).Serve(server.managementPort)
	go grpc.New(server.logger, server.crypto, server.token).Serve(server.grpcPort)
}

func (server *Server) shutdown() {}
