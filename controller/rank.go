package controller

import (
	"clock/common"
	"clock/model"
	"clock/vo"
	"github.com/gin-gonic/gin"
	"time"
)

func Rank(c *gin.Context) {
	db := common.GetDB()
	tel, _ := c.Get("user_tel")
	var curUser model.User
	db.Where("telephone = ?", tel).First(&curUser)
	currentDate := time.Now().Format("2006-01-02")
	var users []model.User
	//db.Preload("Durations", "date = ?", currentDate).Order("time_total desc").Find(&users)
	if err := db.Joins("JOIN durations ON users.telephone = durations.tel").
		Where("durations.date = ?", currentDate).
		Order("durations.dur desc").
		Preload("Durations").
		Find(&users).Error; err != nil {
			c.JSON(500, gin.H{
				"error" : "查询出错",
			})
			return
		}
	var simpleUsers []model.SimpleUser
	for _, user := range users {
		simUser := model.SimpleUser{
			Username:  user.Name,
			TimeTotal: user.Durations[0].Dur,
		}
		simpleUsers = append(simpleUsers, simUser)
	}

	rank := vo.Rank{
		RankMsg: simpleUsers,
		Name:    curUser.Name,
		Date:    currentDate,
	}

	c.JSON(200, rank)
}
