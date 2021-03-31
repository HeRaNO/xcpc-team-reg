package util

import (
	"bytes"

	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type R map[string]interface{}

type contextKey struct {
	name string
}

var StatusCtxKey = &contextKey{"Status"}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := jsoniter.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}

func SuccessResponseWithTotal(w http.ResponseWriter, r *http.Request, data interface{}, total int) {
	respData := R{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		respData["data"] = []interface{}{}
	} else {
		respData["data"] = data
		respData["total"] = total
	}
	JSON(w, r, respData)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	respData := R{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		respData["data"] = []interface{}{}
	} else {
		respData["data"] = data
	}
	JSON(w, r, respData)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, msg string, code string) {
	respData := R{
		"code": code,
		"msg":  msg,
	}
	JSON(w, r, respData)
}
