package main

import (
	"github.com/dylanmazurek/supernote-sync/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = logger.New()

	log.Info().Msg("TBD")
}
