package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type SupernoteClaims struct {
	CreatedTime int64  `json:"createTime"`
	ExpiresAt   int64  `json:"exp"`
	UserId      string `json:"userId"`

	jwt.MapClaims
}
