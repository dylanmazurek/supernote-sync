package models

type Meta struct {
	Success *bool `json:"success"`
	Error   *Error
}

type Error struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	ErrorMsg  *string `json:"errorMsg,omitempty"`
}
