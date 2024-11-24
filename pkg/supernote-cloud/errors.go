package supernotecloud

import "errors"

var (
	ErrAccessTokenExpired = errors.New("access token has expired")
)
