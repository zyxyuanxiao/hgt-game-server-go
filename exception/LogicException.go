package exception

import "net/http"

type LogicException struct {
	Message string
	Code int
	Status int
	ReturnType string
}

// 逻辑异常
func Logic(message string) {
	panic(LogicException{
		Message: message,
		Code: 100,
		Status: http.StatusOK,
		ReturnType: "LogicException",
	})
}