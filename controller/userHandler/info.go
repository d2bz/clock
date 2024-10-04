package userHandler

import (
	"clock/model"
	"clock/responseObject"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	curUser, _ := c.Get("curUser")
	var reUser = responseObject.ReUser{
		Name:      curUser.(model.User).Name,
		Telephone: curUser.(model.User).Telephone,
		Avatar:    curUser.(model.User).Avatar,
	}
	c.JSON(200, reUser)
}
