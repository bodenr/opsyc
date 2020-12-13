package util

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	binary := filepath.Base(os.Args[0])
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Log = zerolog.New(os.Stdout).With().Timestamp().Str(
		"binary", binary).Str("hostname", hostname).Logger()
}
