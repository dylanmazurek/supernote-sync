package main

import (
	"context"
	"os"

	"github.com/dylanmazurek/supernote-sync/internal/logger"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote"
	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	log.Logger = logger.New()

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	supernoteClient, err := supernote.New(ctx, supernote.WithCredentials(username, password))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create supernote client")
	}

	err = supernoteClient.DownloadAllFiles("0", ".files")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to download files")
	}

	// tab := tabulate.New(tabulate.Unicode)
	// tab.Header("Key").SetAlign(tabulate.ML)
	// tab.Header("Value")

	// err = tabulate.Reflect(tab, 0, nil, fileList)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to print table")
	// }

	// tab.Print(os.Stdout)
}
