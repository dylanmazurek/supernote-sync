package noteparser

import (
	"fmt"
	"log"

	"github.com/bzick/tokenizer"
	"github.com/dylanmazurek/supernote-sync/pkg/noteparser/models"
)

// define custom tokens keys
const (
	TNoteStart tokenizer.TokenKey = iota + 1
	TColon
	TArrowOpen
	TArrowClose

	TVersionPrefix
	TWhitespace
	TUnderscore
	TControlChars

	TBinary
)

type noteParser struct {
	tokenizer *tokenizer.Tokenizer
}

var (
	NEW_LINE_CONTROL_CHAR      = "\x01\x00\x00"
	END_OF_LAYERS_CONTROL_CHAR = "\x7f\x00\x00\x00"
)

// NewNoteParser create and configure new tokenizer for Note.
func NewNoteParser() *noteParser {
	parser := &noteParser{}
	parser.tokenizer = tokenizer.New()
	parser.tokenizer.
		DefineTokens(TBinary, []string{"b\xef", "?", "\xef", "\xbf", "\xbd"}).
		DefineTokens(TColon, []string{":"}).
		DefineTokens(TArrowOpen, []string{"<"}).
		DefineTokens(TArrowClose, []string{">"}).
		DefineTokens(TNoteStart, []string{"note"}).
		DefineTokens(TVersionPrefix, []string{"SN_FILE_VER_"}).
		DefineTokens(TControlChars, []string{
			NEW_LINE_CONTROL_CHAR,
			END_OF_LAYERS_CONTROL_CHAR,
		})

	parser.tokenizer.AllowKeywordUnderscore()

	return parser
}

func (parser *noteParser) Parse(noteBytes []byte) (interface{}, error) {
	parsedBytes := parser.tokenizer.ParseBytes(noteBytes)
	parsedNote, err := parser.analyzer(parsedBytes)
	if err != nil {
		return nil, err
	}

	return parsedNote, nil
}

func (parser *noteParser) analyzer(stream *tokenizer.Stream) (*models.Note, error) {
	var newNote models.Note

	for {
		if stream.CurrentToken().Is(TControlChars) {
			log.Printf("TControlChars")
			stream.GoNext()
		} else if stream.CurrentToken().Is(TNoteStart) {
			log.Printf("TNoteStart")
			parseHeader(&newNote, stream)
		} else if stream.CurrentToken().Is(TArrowOpen) {
			log.Printf("TArrowOpen")
			parseElements(&newNote, stream)
		} else if stream.CurrentToken().Is(TBinary, tokenizer.TokenKeyword) {
			log.Printf("TBinary")
			parseBytes(&newNote, stream)
		} else {
			log.Printf("Unknown")
			break
		}
	}

	return nil, parser.error(stream)
}

func parseHeader(note *models.Note, stream *tokenizer.Stream) error {
	stream.GoNext()

	for {
		if stream.CurrentToken().Is(TVersionPrefix) {
			log.Printf("TVersionPrefix")
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenInteger) {
			log.Printf("TokenInteger")
			note.FileVersion = stream.CurrentToken().ValueString()
			stream.GoNext()
		} else if stream.CurrentToken().Is(TControlChars) {
			log.Printf("TControlChars")
			break
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			log.Printf("TokenUndef")
			break
		}
	}

	stream.GoNext()

	return nil
}

func parseElements(note *models.Note, stream *tokenizer.Stream) error {
	stream.GoNext()

	var elements map[string]string = make(map[string]string)

	var currentKey, currentVal string
	for {
		if stream.CurrentToken().Is(TArrowOpen) {
			log.Printf("TArrowOpen")
			stream.GoNext()
		} else if stream.CurrentToken().Is(TArrowClose) {
			log.Printf("TArrowClose")
			elements[currentKey] = currentVal
			currentKey = ""
			currentVal = ""
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenKeyword, tokenizer.TokenInteger) {
			log.Printf("TokenKeyword")
			if currentKey == "" {
				log.Printf("ElementKey")
				currentKey = stream.CurrentToken().ValueString()
			} else {
				log.Printf("ElementValue")
				currentVal = stream.CurrentToken().ValueString()
			}
			stream.GoNext()
		} else if stream.CurrentToken().Is(TColon) {
			log.Printf("TColon")
			stream.GoNext()
		} else if stream.CurrentToken().Is(TControlChars) {
			log.Printf("TControlChars")
			break
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			log.Printf("TokenUndef")
			break
		}
	}

	note.Elements = elements

	stream.GoNext()

	return nil
}

func parseBytes(note *models.Note, stream *tokenizer.Stream) error {
	stream.GoNext()

	var layers []models.Layer

	var currentLayer models.Layer
	for {
		if stream.CurrentToken().Is(TControlChars) {
			log.Printf("TControlChars")

			layers = append(layers, currentLayer)
			currentLayer = models.Layer{}

			if string(stream.CurrentToken().Value()) == END_OF_LAYERS_CONTROL_CHAR {
				log.Printf("END_OF_LAYERS_CONTROL_CHAR")
				break
			}

			stream.GoNext()
			continue
		} else if stream.CurrentToken().Is(TBinary, tokenizer.TokenKeyword) {
			log.Printf("TBinary")
			currentLayer.Bytes = append(currentLayer.Bytes, stream.CurrentToken().Value()...)
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			log.Printf("TokenUndef")
			break
		}
	}

	note.Layers = layers

	stream.GoNext()

	return nil
}

func (parser *noteParser) error(stream *tokenizer.Stream) error {
	isValid := stream.IsValid()
	snippetStr := stream.GetSnippetAsString(5, 0, 0)
	tokenVal := stream.CurrentToken().Value()
	if isValid {
		err := fmt.Errorf("unexpected token %s on line %d near: %s <-- there", tokenVal, stream.CurrentToken().Line(), snippetStr)
		return err
	}

	err := fmt.Errorf("unexpected end on line %d near: %s <-- there", stream.CurrentToken().Line(), snippetStr)
	return err
}
