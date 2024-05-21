package models

import "net/http"

type Request struct {
	HTTPRequest *http.Request
}
