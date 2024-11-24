package main

import (
	"fmt"
	"os"

	"github.com/dylanmazurek/supernote-sync/internal/logger"
	"github.com/dylanmazurek/supernote-sync/pkg/noteparser"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = logger.New()

	storageDir := "/home/dylan/dev/supernote-sync/.storage"
	filePath := fmt.Sprintf(storageDir + "/base_note.note")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("error reading file")
	}

	parser := noteparser.NewNoteParser()
	note, err := parser.Parse(fileContent)
	if err != nil {
		log.Error().Err(err).Msg("error parsing file")
	}

	log.Info().Msgf(fmt.Sprintf("%+v", note))
}
