package main

import (
	"github.com/dylanmazurek/supernote-sync/internal/logger"
	"github.com/dylanmazurek/supernote-sync/pkg/noteparser"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = logger.New()
	noteparser.BomTest()

	// storageDir := "/home/dylan/dev/supernote-sync/.storage"
	// //filePath := fmt.Sprintf(storageDir + "/base_note_2.note")
	// filePath := fmt.Sprintf(storageDir + "/test_no_line.note")

	// fileContent, err := os.ReadFile(filePath)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error reading file")
	// }

	// // for _, fileByte := range fileContent {
	// // 	encodedString := hex.EncodeToString([]byte{fileByte})

	// // 	fmt.Printf("%s\n", encodedString)
	// // }

	// parser := noteparser.NewNoteParser()
	// note, err := parser.Parse(fileContent)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error parsing file")
	// }

	// log.Info().Msgf(fmt.Sprintf("%+v", note))
}
