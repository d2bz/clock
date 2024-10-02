package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReJson 定义返回结构体
type ReJson struct {
	Code  int `json:"code"`  //返回代码
	Msg   any `json:"msg"`   //返回提示信息
	Count int `json:"count"` //返回总数
	Data  any `json:"data"`  //返回数据
}

// Success 请求成功的返回体，传入请求成功的数据和总数
func Success(c *gin.Context, data any, count int) {
	//将参数赋值给你结构体
	Json := ReJson{
		Code:  200,
		Msg:   "请求成功！",
		Count: count,
		Data:  data,
	}
	c.JSON(http.StatusOK, Json)
}

func ResOk(c *gin.Context) {
	//将参数赋值给你结构体
	Json := ReJson{
		Code:  200,
		Msg:   "请求成功！",
		Count: 0,
		Data:  nil,
	}
	c.JSON(http.StatusOK, Json)
}

// Error 请求成功但是有错误的返回体，把错误提示信息传入就行
func Error(c *gin.Context, msg any) {
	//将参数赋值给你结构体
	Json := ReJson{
		Code:  400,
		Msg:   msg,
		Count: 0,
		Data:  nil,
	}
	c.JSON(http.StatusOK, Json)
}

// ResErr Fail 请求失败的返回体（网络不通），只需要传入请求失败的信息回来就行了
func ResErr(c *gin.Context) {
	//将参数赋值给你结构体
	Json := ReJson{
		Code:  400,
		Msg:   "请求失败！",
		Count: 0,
		Data:  nil,
	}
	c.JSON(http.StatusNotFound, Json)

}

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
