package main

import (
	"crypto/md5"
	"os"
	"unicode/utf8"

	"github.com/dylanmazurek/supernote-sync/internal/logger"
	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog/log"
)

func main() {
	//ctx := context.Background()
	log.Logger = logger.New()

	log.Info().Msg("starting supernote-sync")

	//username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	hashedPassword := testPassword(password)
	log.Info().Msg(hashedPassword)

	// _, err := supernote.New(ctx, supernote.WithCredentials(username, password))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to create supernote client")
	// }

	// log.Info().Msg("successfully authenticated")
}

func testPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	hashPassword := string(hash.Sum(nil))
	log.Info().Msg(hashPassword)

	r, _ := utf8.DecodeRuneInString(hashPassword)
	log.Info().Msg(string(r))

	return hashPassword
}
