package controller

import (
	"clock/common"
	"clock/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Start(c *gin.Context) {
	tel, _ := c.Get("user_tel")
	db := common.GetDB()
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")
	// 查询当前用户当天是否有打卡记录，有则更新，无则新增
	var duration model.Duration
	db.Where("telephone = ? AND date = ?", tel, currentDate).First(&duration)
	// 开始时间不为空，说明正在计时中
	if duration.StartTime != "" {
		common.Error(c, "已完成开始打卡")
		return
	}
	// 开始时间为空则判断今日是否打过卡，打过卡说明已经创建过记录，把记录更新即可
	if duration.ID != 0 {
		db.Model(&duration).Update("start_time", currentTime)
		common.Success(c, "开始打卡成功", 1)
		return
	}
	// 未打过卡就新增本日记录
	newDur := model.Duration{
		Tel:       tel.(string),
		Date:      currentDate,
		StartTime: currentTime,
		Dur:       0,
	}
	db.Create(&newDur)
	common.Success(c, "开始打卡成功", 1)
}

func End(c *gin.Context) {
	tel, _ := c.Get("user_tel")
	db := common.GetDB()
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")
	// 查询当前用户当天是否有打卡记录
	var duration model.Duration
	db.Where("telephone = ? AND date = ?", tel, currentDate).First(&duration)
	// 提示未进行开始打卡
	if duration.StartTime == "" {
		common.Error(c, "未进行开始打卡")
		return
	}
	// 结算本次打卡时间
	t1, err1 := time.Parse("15:04:05", duration.StartTime)
	t2, err2 := time.Parse("15:04:05", currentTime)
	if err1 != nil || err2 != nil {
		common.Error(c, "时间解析出错")
		return
	}
	durObj := t2.Sub(t1)
	dur := int(durObj.Seconds())
	todayDur := duration.Dur + dur
	// 更新本日记录
	db.Model(&duration).Updates(map[string]interface{}{
		"StartTime": "",
		"Dur":       todayDur,
	})
	// 返回具体时分秒
	h := dur / 3600
	m := (dur % 3600) / 60
	//s := dur % 60
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "结束打卡成功",
		"hours":   h,
		"minutes": m,
		//"seconds": s,
	})
}
