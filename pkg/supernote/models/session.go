package models

import (
	"github.com/dylanmazurek/supernote-sync/pkg/utilities/passwordhash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type Session struct {
	username string
	password string

	randomCode string
	timestamp  int64

	token jwt.Token
}

func (s *Session) SetCredentials(username, password string) {
	s.username = username
	s.password = password
}

func (s *Session) GetToken() string {
	return s.token.Raw
}

func (s *Session) SetToken(token string) error {
	validator := jwt.NewParser(jwt.WithoutClaimsValidation(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	claims := jwt.MapClaims{}
	parsedToken, _, err := validator.ParseUnverified(token, claims)
	if err != nil {
		return err
	}

	s.token = *parsedToken

	expiry, err := s.token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}

	log.Debug().Time("expires_at", expiry.Time).Msg("token set successfully")

	return nil
}

func (s *Session) SetMetadata(randomCode string, timestamp int64) {
	s.randomCode = randomCode
	s.timestamp = timestamp
}

func (s *Session) GetMetadata() (string, int64) {
	hashedPassword := passwordhash.HashPassword(s.randomCode, s.password)

	return hashedPassword, s.timestamp
}
