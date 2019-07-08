package sealog_errors

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var errorMsgs = map[string]errorMsg{
	"default_error": errorMsg{
		Code:    1001,
		Message: "default_error",
	},

	"params_error": errorMsg{
		Code:    1002,
		Message: "参数解析有误",
	},
}

func GetError(err error) errorMsg {
	if err == nil {
		return errorMsgs["default_error"]
	}
	return errorMsgs[err.Error()]
}
