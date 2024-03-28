package utils

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

type R map[string]interface{}

func ErrorResp(err berrors.Berror) R {
	return R{
		"code": err.Code(),
		"msg":  err.Msg(),
	}
}

func SuccessResp(data interface{}) R {
	resp := R{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		resp["data"] = []interface{}{}
	} else {
		resp["data"] = data
	}
	return resp
}
