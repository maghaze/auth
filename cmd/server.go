package cmd

import (
	"os"

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
	managementPort int
	grpcPort       int
}

func (server Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		server.main(config.Load(true), trap)
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "serve the authentication and authorization server",
		Run:   run,
	}

	cmd.Flags().IntVar(&server.managementPort, "management-port", 8080, "The port the metrics and probe endpoint binds to")
	cmd.Flags().IntVar(&server.grpcPort, "grpc-port", 9090, "The port the grpc endpoint listens on")

	return cmd
}

func (server *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	crypto := crypto.New(cfg.Crypto)
	token, err := token.New(cfg.Token)
	if err != nil {
		logger.Panic("Error creating token object", zap.Error(err))
	}

	go rest.New(logger).Serve(server.managementPort)
	go grpc.New(logger, crypto, token).Serve(server.grpcPort)

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)
}
