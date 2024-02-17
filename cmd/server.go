package cmd

import (
	"os"

	"github.com/maghaze/auth/internal/config"
	"github.com/maghaze/auth/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct{}

func (cmd Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.main(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "serve the service",
		Run:   run,
	}
}

func (cmd *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	// crypto := crypto.New(cfg.Crypto)
	// token, err := token.New(cfg.Token)
	// if err != nil {
	// 	logger.Panic("Error creating token object", zap.Error(err))
	// }

	// go grpc.New(cfg.Grpc, logger, crypto, token).Serve()

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)
}
