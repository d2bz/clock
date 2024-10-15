package userHandler

import (
	"clock/Redis"
	"clock/common"
	"clock/model"
	"clock/util"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var telephonePattern = `^\d{11}$`
var passwordPattern = `^[a-zA-Z0-9!@#$%^&*()_+={}|[\]\\:";'<>?,./]{6,}$`

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name") //注意字符串要用双引号
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if matched, err := isValidForm(telephonePattern, telephone); !matched {
		if err == nil {
			util.Response(ctx, http.StatusBadRequest, "手机号必须为11位数字", "")
			return
		} else {
			util.Response(ctx, http.StatusInternalServerError, "手机号格式匹配出错", err.Error())
			return
		}
	}

	if matched, err := isValidForm(passwordPattern, password); !matched {
		if err == nil {
			util.Response(ctx, http.StatusBadRequest, "密码格式错误", "")
			return
		} else {
			util.Response(ctx, http.StatusInternalServerError, "密码格式匹配出错", err.Error())
			return
		}
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		util.Response(ctx, http.StatusBadRequest, "手机号已被注册", "")
		return
	}
	//如果用户不存在，则创建用户
	//对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.Response(ctx, http.StatusInternalServerError, "密码加密错误", err.Error())
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
	err = Redis.SetUserInfo(newUser)
	if err != nil {
		util.Response(ctx, http.StatusInternalServerError, "redis userinfo set出错", err.Error())
		return
	}
	util.Response(ctx, http.StatusOK, "注册成功", "")
}

// Login 登录功能
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	tel := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if matched, err := isValidForm(telephonePattern, tel); !matched {
		if err == nil {
			util.Response(ctx, http.StatusBadRequest, "手机号必须为11位数字", "")
			return
		} else {
			util.Response(ctx, http.StatusInternalServerError, "手机号格式匹配出错", err.Error())
			return
		}
	}

	if matched, err := isValidForm(passwordPattern, password); !matched {
		if err == nil {
			util.Response(ctx, http.StatusBadRequest, "密码格式错误", "")
			return
		} else {
			util.Response(ctx, http.StatusInternalServerError, "密码格式匹配出错", err.Error())
			return
		}
	}

	//判断用户是否存在
	var user model.User
	temp, err := Redis.GetUserInfo(tel)
	if err == nil {
		user = temp
	} else {
		DB.Where("telephone = ?", tel).First(&user)

		if user.ID == 0 {
			util.Response(ctx, http.StatusBadRequest, "用户不存在", err.Error())
			return
		}
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		util.Response(ctx, http.StatusBadRequest, "密码错误", "")
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

// 判断格式是否合法
func isValidForm(pattern string, s string) (bool, error) {

	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0
}
