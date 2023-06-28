package utils

type R map[string]interface{}

func ErrorResp(errcode string, msg string) R {
	return R{
		"code": errcode,
		"msg":  msg,
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
