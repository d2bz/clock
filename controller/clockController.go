package controller

import (
	"clock/Redis"
	"clock/common"
	"clock/model"
	"clock/util"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type startJson struct {
	Position string `json:"position" binding:"required"`
}

func Start(c *gin.Context) {
	curUser, _ := c.Get("curUser")
	uid := curUser.(model.User).UserID
	db := common.GetDB()
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")
	//判断是否已经进行开始打卡
	//redis初步筛选
	if flag, _ := Redis.GetIsClock(uid); flag {
		util.Response(c, http.StatusBadRequest, "正在计时中", "")
		return
	}

	var duration model.Duration
	db.Where("user_id = ? AND date = ? AND end_time = ?", uid, currentDate, "").First(&duration)
	if duration.ID != 0 {
		util.Response(c, http.StatusBadRequest, "正在计时中", "")
		return
	}

	//绑定前端传入信息
	var sJson startJson
	if err := c.ShouldBindJSON(&sJson); err != nil {
		util.Response(c, http.StatusInternalServerError, "接收Json出错", err.Error())
		return
	}
	newDur := model.Duration{
		UserID:    uid,
		Date:      currentDate,
		StartTime: currentTime,
		EndTime:   "",
		Dur:       0,
		Position:  sJson.Position,
	}

	if err := Redis.SetIsClock(uid); err != nil {
		util.Response(c, http.StatusInternalServerError, "redis isClock set失败", err.Error())
		return
	}
	db.Create(&newDur)
	util.Response(c, http.StatusOK, "开始打卡成功", "")

}

func End(c *gin.Context) {
	curUser, _ := c.Get("curUser")
	uid := curUser.(model.User).UserID
	db := common.GetDB()
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")
	// 查询当前用户当天是否有打卡记录
	//redis初步筛选
	flag, err := Redis.GetIsClock(uid)
	if !flag && err == redis.Nil {
		util.Response(c, http.StatusBadRequest, "未进行开始打卡", "redis返回")
		return
	}

	//排查redis是否出现其他错误
	if err != nil && err != redis.Nil {
		log.Printf("err: %v", err.Error())
	}

	var duration model.Duration
	db.Where("user_id = ? AND date = ? AND end_time = ?", uid, currentDate, "").First(&duration)
	// 提示未进行开始打卡
	if duration.ID == 0 {
		util.Response(c, http.StatusBadRequest, "未进行开始打卡", "")
		return
	}
	// 结算本次打卡时间
	startTime, err1 := time.Parse("15:04:05", duration.StartTime)
	endTime, err2 := time.Parse("15:04:05", currentTime)
	if err1 != nil || err2 != nil {
		util.Response(c, http.StatusInternalServerError, "时间解析出错", "")
		return
	}
	durObj := endTime.Sub(startTime)
	dur := int(durObj.Seconds())

	// 转换成具体分钟(min)
	m := dur / 60

	//删除redis的键
	err = Redis.DeleteIsClock(uid)
	if err != nil {
		util.Response(c, http.StatusInternalServerError, "redis isClock删除失败", err.Error())
	}

	// 更新本日记录
	db.Model(&duration).Updates(map[string]interface{}{
		"EndTime": currentTime,
		"Dur":     m,
	})

	mString := strconv.Itoa(m)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "结束打卡成功",
		"minutes": mString + "min",
	})

}

// func Start(c *gin.Context) {
// 	curUser, _ := c.Get("curUser")
// 	tel := curUser.(model.User).Telephone
// 	db := common.GetDB()
// 	currentDate := time.Now().Format("2006-01-02")
// 	currentTime := time.Now().Format("15:04:05")
// 	// 查询当前用户当天是否有打卡记录，有则更新，无则新增
// 	var duration model.Duration
// 	db.Where("tel = ? AND date = ?", tel, currentDate).First(&duration)
// 	// 开始时间不为空，说明正在计时中
// 	if duration.StartTime != "" {
// 		common.Error(c, "已完成开始打卡")
// 		return
// 	}
// 	// 开始时间为空则判断今日是否打过卡，打过卡说明已经创建过记录，把记录更新即可
// 	if duration.ID != 0 {
// 		db.Model(&duration).Update("start_time", currentTime)
// 		common.Success(c, "开始打卡成功", 1)
// 		return
// 	}
// 	// 未打过卡就新增本日记录
// 	newDur := model.Duration{
// 		Tel:       tel,
// 		Date:      currentDate,
// 		StartTime: currentTime,
// 		Dur:       0,
// 	}
// 	db.Create(&newDur)
// 	common.Success(c, "开始打卡成功", 2)
// }

// func End(c *gin.Context) {
// 	curUser, _ := c.Get("curUser")
// 	tel := curUser.(model.User).Telephone
// 	db := common.GetDB()
// 	currentDate := time.Now().Format("2006-01-02")
// 	currentTime := time.Now().Format("15:04:05")
// 	// 查询当前用户当天是否有打卡记录
// 	var duration model.Duration
// 	db.Where("tel = ? AND date = ?", tel, currentDate).First(&duration)
// 	// 提示未进行开始打卡
// 	if duration.StartTime == "" {
// 		common.Error(c, "未进行开始打卡")
// 		return
// 	}
// 	// 结算本次打卡时间
// 	t1, err1 := time.Parse("15:04:05", duration.StartTime)
// 	t2, err2 := time.Parse("15:04:05", currentTime)
// 	if err1 != nil || err2 != nil {
// 		common.Error(c, "时间解析出错")
// 		return
// 	}
// 	durObj := t2.Sub(t1)
// 	dur := int(durObj.Seconds())

// 	// 转换成具体分钟(min)
// 	m := dur / 60

// 	todayDur := duration.Dur + m
// 	// 更新本日记录
// 	db.Model(&duration).Updates(map[string]interface{}{
// 		"StartTime": "",
// 		"Dur":       todayDur,
// 	})

// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    200,
// 		"msg":     "结束打卡成功",
// 		"minutes": m,
// 	})
// }
