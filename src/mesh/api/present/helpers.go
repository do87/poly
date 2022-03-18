package present

import (
	"encoding/json"
	"net/http"
)

// Presentor is the response struct
type SupportedDataTypes interface {
	interface{} | string | TypeRun | []TypeRun | TypeKey | []TypeKey | TypeAgent | []TypeAgent
}
type Presentor[K SupportedDataTypes] struct {
	Kind string `json:"kind"`
	Etag string `json:"etag,omitempty"`
	Data K      `json:"data"`
}

func wrap[K SupportedDataTypes](kind, etag string, data K) Presentor[K] {
	return Presentor[K]{
		Kind: kind,
		Etag: etag,
		Data: data,
	}
}

// Generic makes wrap functionality public
var Generic = wrap[interface{}]

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

func Unmarshal[K SupportedDataTypes](b []byte) (p Presentor[K], err error) {
	err = json.Unmarshal(b, &p)
	return
}

func GetDataFromBytes[K SupportedDataTypes](b []byte) (d K, err error) {
	p, err := Unmarshal[K](b)
	if err != nil {
		return
	}
	return p.Data, nil
}
