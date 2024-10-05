package userHandler

import (
	"clock/common"
	"clock/model"
	"clock/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ModifyPassword(c *gin.Context) {
	//curUser, _ := c.Get("curUser")

}

func ModifyShowInfo(c *gin.Context) {
	curUser, _ := c.Get("curUser")
	newName := c.PostForm("newName")
	//newTelephone := c.PostForm("newTelephone")

	db := common.GetDB()
	var user model.User
	db.Where("telephone = ?", curUser.(model.User).Telephone).First(&user)

	err := db.Model(&user).Updates(map[string]interface{}{
		"name": newName,
	}).Error
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "更新失败", "")
	}

	util.Response(c, http.StatusOK, "更新成功", "")
}
