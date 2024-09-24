package controller

import (
	"clock/common"
	"clock/model"
	"clock/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name") //注意字符串要用双引号
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		common.Error(ctx, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		common.Error(ctx, "密码不能少于6位")
		return
	}
	//如果没有传name（名称），就给他一个10位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		common.Error(ctx, "用户已经存在")
		return
	}
	//如果用户不存在，则创建用户
	//对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		common.Error(ctx, "加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword), //将加密后的密码保存起来
	}
	DB.Create(&newUser)
	//返回结果
	//ctx.JSON(200, gin.H{
	//	"message": "注册成功",
	//})
	common.Success(ctx, "注册成功", 1)
}

// Login 登录功能
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	if len(password) < 6 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		common.Error(ctx, "密码不能少于6位")
		return
	}
	//判断用户是否存在
	var user model.User
	DB.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		common.Error(ctx, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		common.Error(ctx, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
	//common.Success(ctx, "登录成功", 1)
	//tel := user.Telephone
	//setCookie(tel, ctx)
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0
}

// func setCookie(tel string, c *gin.Context) {
// 	cookie := &http.Cookie{
// 		Name:     "Telephone",
// 		Value:    tel,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(c.Writer, cookie)
// 	c.JSON(http.StatusOK, gin.H{
// 		"msg": "Cookie已设置",
// 	})
// }
