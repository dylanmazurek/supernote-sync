package noteparser

import (
	"fmt"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/bzick/tokenizer"
	"github.com/dylanmazurek/supernote-sync/pkg/noteparser/models"
)

// define custom tokens keys
const (
	TNoteStart tokenizer.TokenKey = iota + 1
	TColon
	TArrowOpen
	TArrowClose
	TCurlyOpen
	TCurlyClose
	TSquareOpen
	TSquareClose
	TComma
	TDoubleQuote

	TVersionPrefix
	TWhitespace
	TUnderscore
	TControlChars

	TBinary
)

var tokenKeyNames = map[tokenizer.TokenKey]string{
	tokenizer.TokenUnknown:        "TokenUnknown",
	tokenizer.TokenStringFragment: "TokenStringFragment",
	tokenizer.TokenString:         "TokenString",
	tokenizer.TokenFloat:          "TokenFloat",
	tokenizer.TokenInteger:        "TokenInteger",
	tokenizer.TokenKeyword:        "TokenKeyword",
	tokenizer.TokenUndef:          "TokenUndef",

	TNoteStart:   "TNoteStart",
	TColon:       "TColon",
	TArrowOpen:   "TArrowOpen",
	TArrowClose:  "TArrowClose",
	TCurlyOpen:   "TCurlyOpen",
	TCurlyClose:  "TCurlyClose",
	TSquareOpen:  "TSquareOpen",
	TSquareClose: "TSquareClose",
	TComma:       "TComma",
	TDoubleQuote: "TDoubleQuote",

	TVersionPrefix: "TVersionPrefix",
	TWhitespace:    "TWhitespace",
	TUnderscore:    "TUnderscore",
	TControlChars:  "TControlChars",

	TBinary: "TBinary",
}

type noteParser struct {
	tokenizer *tokenizer.Tokenizer
}

var (
	CONTROL_CHAR_NEW_LINE            = "\x01\x00\x00"
	CONTROL_CHAR_END_OF_LAYERS       = "\x7f\x00\x00\x00"
	CONTROL_CHAR_END_OF_LAYER_META   = "}\x00\x00\x00"
	CONTROL_CHAR_END_OF_LAYER_LIST   = "\xef\xbf"
	CONTROL_CHAR_END_OF_LAYER_LIST_2 = "\xef\xdb"
	CONTROL_CHAR_END_OF_LAYER_LIST_3 = "\xbd\x04\x00\x00"
	CONTROL_CHAR_END_OF_LAYER_LIST_4 = "\xef\xbf\x04\x00\x00"
	CONTROL_CHAR_TAIL_OF_FILE        = "=\x00\x00\x00"
	CONTROL_CHAR_TAIL_OF_FILE_2      = "\xef\xbf\xbd"
)

// NewNoteParser create and configure new tokenizer for Note.
func NewNoteParser() *noteParser {
	parser := &noteParser{}
	parser.tokenizer = tokenizer.New()
	parser.tokenizer.
		DefineTokens(TControlChars, []string{
			CONTROL_CHAR_NEW_LINE,
			CONTROL_CHAR_END_OF_LAYERS,
			CONTROL_CHAR_END_OF_LAYER_META,
			CONTROL_CHAR_END_OF_LAYER_LIST,
			CONTROL_CHAR_END_OF_LAYER_LIST_2,
			CONTROL_CHAR_END_OF_LAYER_LIST_3,
			CONTROL_CHAR_END_OF_LAYER_LIST_4,
			CONTROL_CHAR_TAIL_OF_FILE,
			CONTROL_CHAR_TAIL_OF_FILE_2,
		}).
		DefineTokens(TBinary, []string{"b\xef", "?", "\xbf", "\xbd"}).
		DefineTokens(TColon, []string{":"}).
		DefineTokens(TArrowOpen, []string{"<"}).
		DefineTokens(TArrowClose, []string{">"}).
		DefineTokens(TCurlyOpen, []string{"{"}).
		DefineTokens(TCurlyClose, []string{"}"}).
		DefineTokens(TSquareOpen, []string{"["}).
		DefineTokens(TSquareClose, []string{"]"}).
		DefineTokens(TComma, []string{","}).
		DefineTokens(TNoteStart, []string{"note"}).
		DefineTokens(TVersionPrefix, []string{"SN_FILE_VER_"})

	parser.tokenizer.AllowKeywordUnderscore()
	parser.tokenizer.DefineStringToken(TDoubleQuote, `"`, `"`).SetEscapeSymbol(tokenizer.BackSlash)

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
			printDebug(stream.CurrentToken(), nil)
			stream.GoNext()
		} else if stream.CurrentToken().Is(TNoteStart) {
			printDebug(stream.CurrentToken(), nil)
			parseHeader(&newNote, stream)
		} else if stream.CurrentToken().Is(TArrowOpen) {
			printDebug(stream.CurrentToken(), nil)
			parseElements(&newNote, stream)
		} else if stream.CurrentToken().Is(TBinary, tokenizer.TokenKeyword) {
			printDebug(stream.CurrentToken(), nil)
			parseBytes(&newNote, stream)
		} else {
			printDebug(stream.CurrentToken(), nil)
			break
		}
	}

	return nil, parser.error(stream)
}

func parseHeader(note *models.Note, stream *tokenizer.Stream) error {
	stream.GoNext()

	for {
		if stream.CurrentToken().Is(TVersionPrefix) {
			printDebug(stream.CurrentToken(), nil)
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenInteger) {
			printDebug(stream.CurrentToken(), nil)
			note.FileVersion = stream.CurrentToken().ValueString()
			stream.GoNext()
		} else if stream.CurrentToken().Is(TControlChars) {
			printDebug(stream.CurrentToken(), nil)
			break
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			printDebug(stream.CurrentToken(), ptr.String("parse header"))
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
		if stream.CurrentToken().Is(TArrowOpen, TArrowClose) {
			if stream.CurrentToken().Is(TArrowOpen) {
				printDebug(stream.CurrentToken(), nil)
			} else {
				printDebug(stream.CurrentToken(), nil)
				elements[currentKey] = currentVal
				currentKey = ""
				currentVal = ""
			}

			stream.GoNext()
			// } else if stream.CurrentToken().Is(TCurlyOpen, TCurlyClose) {
			// 	if stream.CurrentToken().Is(TCurlyOpen) {
			// 		printDebug(stream.CurrentToken())
			// 	} else {
			// 		printDebug(stream.CurrentToken())
			// 	}

			// 	stream.GoNext()
		} else if stream.CurrentToken().Is(TSquareOpen, TSquareClose) {
			parseJsonModified(note, stream)
			// if stream.CurrentToken().Is(TSquareOpen) {
			// 	printDebug(stream.CurrentToken())
			// } else {
			// 	printDebug(stream.CurrentToken())
			// }

			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenKeyword, tokenizer.TokenInteger, TComma) {
			if currentKey == "" {
				printDebug(stream.CurrentToken(), nil)
				currentKey = stream.CurrentToken().ValueString()
			} else {
				printDebug(stream.CurrentToken(), nil)
				currentVal = stream.CurrentToken().ValueString()
			}

			stream.GoNext()
		} else if stream.CurrentToken().Is(TColon) {
			printDebug(stream.CurrentToken(), nil)
			stream.GoNext()
		} else if stream.CurrentToken().Is(TControlChars) {
			printDebug(stream.CurrentToken(), nil)
			break
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			printDebug(stream.CurrentToken(), ptr.String("parse elements"))
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
			printDebug(stream.CurrentToken(), nil)

			layers = append(layers, currentLayer)
			currentLayer = models.Layer{}

			if string(stream.CurrentToken().Value()) == CONTROL_CHAR_END_OF_LAYERS {
				printDebug(stream.CurrentToken(), nil)
				break
			}

			stream.GoNext()
			continue
		} else if stream.CurrentToken().Is(TBinary, tokenizer.TokenKeyword) {
			printDebug(stream.CurrentToken(), nil)
			currentLayer.Bytes = append(currentLayer.Bytes, stream.CurrentToken().Value()...)
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			printDebug(stream.CurrentToken(), ptr.String("parse bytes"))
			break
		}
	}

	note.Layers = layers

	stream.GoNext()

	return nil
}

func parseJsonModified(note *models.Note, stream *tokenizer.Stream) error {
	stream.GoNext()

	for {
		if stream.CurrentToken().Is(TSquareClose) {
			printDebug(stream.CurrentToken(), nil)
			stream.GoNext()
		} else if stream.CurrentToken().Is(TCurlyOpen, TCurlyClose) {
			printDebug(stream.CurrentToken(), nil)
			stream.GoNext()
		} else if stream.CurrentToken().Is(tokenizer.TokenUndef) {
			printDebug(stream.CurrentToken(), nil)
			break
		} else {
			printDebug(stream.CurrentToken(), ptr.String("parse json modified"))
			stream.GoNext()
		}
	}

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

func printDebug(token *tokenizer.Token, msg *string) {
	tokenName, ok := tokenKeyNames[token.Key()]
	if !ok {
		tokenName = "Unknown TokenKey"
	}

	var logMsg string
	if msg != nil {
		logMsg = fmt.Sprintf("[%s]", *msg)
	}

	log.Printf("%15s = %-24s%s", tokenName, string(token.Value()), logMsg)
}
