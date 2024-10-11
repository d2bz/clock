package userHandler

import (
	"clock/common"
	"clock/model"
	"clock/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		common.Error(ctx, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		common.Error(ctx, "密码不能少于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		common.Error(ctx, "用户已经存在")
		return
	}
	//如果用户不存在，则创建用户
	//对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		common.Error(ctx, "加密错误")
		return
	}

	userid, err := uuid.NewUUID()
	if err != nil {
		util.Response(ctx, http.StatusInternalServerError, "id生成错误", "")
		return
	}

	newUser := model.User{
		UserID:    userid.String(),
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword), //将加密后的密码保存起来
		Avatar:    "",
	}
	DB.Create(&newUser)

	common.Success(ctx, "注册成功", 1)
}

// Login 登录功能
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	tel := ctx.PostForm("telphone")
	password := ctx.PostForm("password")

	if len(password) < 6 {
		common.Error(ctx, "密码不能少于6位")
		return
	}
	//判断用户是否存在
	var user model.User
	DB.Where("telephone = ?", tel).First(&user)
	if user.ID == 0 {
		common.Error(ctx, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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
