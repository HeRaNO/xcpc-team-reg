package utils

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

type R map[string]any

func ErrorResp(err berrors.Berror) R {
	return R{
		"code": err.Code(),
		"msg":  err.Msg(),
	}
}

func SuccessResp(data any) R {
	resp := R{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		resp["data"] = []any{}
	} else {
		resp["data"] = data
	}
	return resp
}
