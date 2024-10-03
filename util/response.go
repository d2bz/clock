package util

import "github.com/gin-gonic/gin"

type Re_Json struct {
	Code int
	Msg  string
	Data any
}

func Response(c *gin.Context, code int, msg string, data any) {
	Json := Re_Json{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(code, Json)
}
