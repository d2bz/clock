package userHandler

import (
	"clock/common"
	"clock/model"
	"clock/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModifiableUser struct {
	Name      *string
	Telephone *string
}

// type password struct {
// 	OldPassword string
// 	NewPassword string
// 	Confirm     string
// }

func ModifyPassword(c *gin.Context) {
	//curUser, _ := c.Get("curUser")

}

func ModifyShowInfo(c *gin.Context) {
	curUser, _ := c.Get("curUser")
	var form ModifiableUser
	err := c.ShouldBind(&form)
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "数据绑定失败", err.Error())
		return
	}

	// ToDo 从Redis中拿用户信息

	db := common.GetDB()
	var user model.User
	db.Where("user_id = ?", curUser.(model.User).UserID).First(&user)

	updateInfo := map[string]interface{}{}
	if *form.Name != user.Name && form.Name != nil {
		updateInfo["name"] = *form.Name
	}
	if *form.Telephone != user.Telephone && form.Telephone != nil {
		updateInfo["telephone"] = *form.Telephone
	}

	err = db.Model(&user).Updates(updateInfo).Error
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "更新失败", err.Error())
	}

	util.Response(c, http.StatusOK, "更新成功", "")
}

func ModifiAvatar(c *gin.Context) {
	url, err := common.UploadFile(c)
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "文件上传失败", err.Error())
		return
	}

	curUser, _ := c.Get("curUser")

	// ToDo 从Redis中拿用户信息

	db := common.GetDB()
	var user model.User
	db.Where("user_id = ?", curUser.(model.User).UserID).First(&user)
	err = db.Model(&user).Update("avatar", url).Error
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	util.Response(c, http.StatusOK, "更新成功", gin.H{
		"newAvatarUrl": url,
	})
}
