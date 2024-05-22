package main

import (
	"context"
	"os"

	"github.com/dylanmazurek/supernote-sync/internal/logger"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote"
	_ "github.com/joho/godotenv/autoload"

	"github.com/markkurossi/tabulate"
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

	fileList, err := supernoteClient.GetFileList(999713992811216897, 1, 20)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user info")
	}

	tab := tabulate.New(tabulate.Unicode)
	tab.Header("Key").SetAlign(tabulate.ML)
	tab.Header("Value")

	err = tabulate.Reflect(tab, 0, nil, fileList)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to print table")
	}

	tab.Print(os.Stdout)
}
