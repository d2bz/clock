package model

import "gorm.io/gorm"

type Duration struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Date      string `gorm:"type:varchar(100);not null"`
	StartTime string `gorm:"type:varchar(100)"`
	// 打卡时长（秒）
	Dur int `gorm:"type:int"`
}
