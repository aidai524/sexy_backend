package sexyerror

import (
	"encoding/json"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type ThirdPartyError struct {
	Provider string      `json:"provider"`
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Payload  interface{} `json:"-"`
}

func (e *ThirdPartyError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
