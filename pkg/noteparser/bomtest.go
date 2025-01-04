package noteparser

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spkg/bom"
)

func BomTest() {
	storageDir := "/home/dylan/dev/supernote-sync/.storage"
	//filePath := fmt.Sprintf(storageDir + "/base_note_2.note")
	filePath := fmt.Sprintf(storageDir + "/test_no_line.note")

	fileContent, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("error reading file")
	}

	// actual := bom.Clean(tc.Input)

	cleanFileContent := bom.NewReader(fileContent)
	actual, err := io.ReadAll(cleanFileContent)
	if err != nil {
		log.Error().Err(err).Msg("error reading file")
	}

	for _, fileByte := range actual {
		fmt.Printf("%s", string(fileByte))
		// encodedString := hex.EncodeToString([]byte{fileByte})

		// fmt.Printf("%s\n", encodedString)
	}
}
