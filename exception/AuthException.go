package exception

import "net/http"

type AuthException struct {
	Message string
	Code int
	Status int
	ReturnType string
}

// 逻辑异常
func Auth(message string) {
	panic(AuthException{
		Message: message,
		Code: 401,
		Status: http.StatusUnauthorized,
		ReturnType: "AuthException",
	})
}