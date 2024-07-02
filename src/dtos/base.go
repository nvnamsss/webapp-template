package dtos

import "net/http"

type Meta struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}

// NewMeta returns a new meta with message.
func NewMeta(code int, messages ...string) Meta {
	var msg = http.StatusText(code)
	if len(messages) > 0 {
		msg = messages[0]
	}
	return Meta{
		Code:    code,
		Message: msg,
	}
}

type PaginationMeta struct {
	Meta
	Total int64 `json:"total"`
}
