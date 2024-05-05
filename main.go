package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/maghaze/auth/cmd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	const description = "Auth microservice helps to authenticate users as a seperate microservice"
	root := &cobra.Command{Short: description}

	root.AddCommand(
		cmd.Server{}.Command(ctx),
	)

	if err := root.Execute(); err != nil {
		log.Fatal("failed to execute the root command", zap.Error(err))
	}
}
