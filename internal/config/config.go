package config

import (
	"github.com/maghaze/auth/pkg/crypto"
	"github.com/maghaze/auth/pkg/logger"
	"github.com/maghaze/auth/pkg/token"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	Token  *token.Config  `koanf:"token"`
	Crypto *crypto.Config `koanf:"crypto"`
}
