package present

import (
	"net/http"
)

type Presentor struct {
	Kind string      `json:"kind"`
	Etag string      `json:"etag,omitempty"`
	Data interface{} `json:"data"`
}

func wrap(kind, etag string, data interface{}) Presentor {
	return Presentor{
		Kind: kind,
		Etag: etag,
		Data: data,
	}
}

type errorPayload struct {
	HTTPStatus int    `json:"httpStatus"`
	Error      string `json:"error,omitempty"`
}

// Error returns an error response
func Error(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error) errorPayload {
	w.WriteHeader(httpStatusCode)
	return errorPayload{
		HTTPStatus: httpStatusCode,
		Error:      err.Error(),
	}
}
