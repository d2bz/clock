// 创建main。go的第一步就是将package的名称改成main，不然运行的时候，会报系统兼容的错误
package main

import (
	"clock/common"
	"clock/controller"
	"clock/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始服务器引擎
	r := gin.Default()

	db := common.InitDB()
	//获取底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	defer sqlDB.Close()

	controller.StartScheduledDeletion(db)

	r = routes.CollectRoute(r)

	panic(r.Run(":8090"))
}
