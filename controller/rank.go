package controller

import (
	"clock/common"
	"clock/model"
	"clock/vo"
	"time"

	"github.com/gin-gonic/gin"
)

func Rank(c *gin.Context) {
	db := common.GetDB()
	curUser, _ := c.Get("curUser")
	currentDate := time.Now().Format("2006-01-02")
	//尝试不使用外键查询
	var users []model.User
	var simpleUsers []model.SimpleUser
	if err := db.Select("users.name as username, SUM(durations.dur) AS time_total").
		Joins("JOIN durations ON users.telephone = durations.tel").
		Where("durations.date = ?", currentDate).
		Group("users.telephone, users.name").
		Order("time_total desc").
		Find(&users).
		Scan(&simpleUsers).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "查询出错",
		})
		return
	}

	rank := vo.Rank{
		RankMsg: simpleUsers,
		Name:    curUser.(model.User).Name,
		Date:    currentDate,
	}

	c.JSON(200, rank)
}

// func Rank(c *gin.Context) {
// 	db := common.GetDB()
// 	curUser, _ := c.Get("curUser")
// 	// var curUser model.User
// 	// db.Where("telephone = ?", tel).First(&curUser)
// 	currentDate := time.Now().Format("2006-01-02")
// 	//尝试不使用外键查询
// 	var users []model.User
// 	var simpleUsers []model.SimpleUser
// 	if err := db.Select("users.name as username, durations.dur as time_total").
// 		Joins("JOIN durations on users.telephone = durations.tel").
// 		Where("durations.date = ?", currentDate).
// 		Order("durations.dur desc").
// 		Find(&users).
// 		Scan(&simpleUsers).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"error": "查询出错",
// 		})
// 		return
// 	}

// 	rank := vo.Rank{
// 		RankMsg: simpleUsers,
// 		Name:    curUser.(model.User).Name,
// 		Date:    currentDate,
// 	}

// 	c.JSON(200, rank)
// }

//原使用外键
/*
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
*/
