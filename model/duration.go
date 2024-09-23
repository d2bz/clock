package model

import "gorm.io/gorm"

type Duration struct {
	gorm.Model
	Tel       string `gorm:"type:varchar(11);not null"`
	Date      string `gorm:"type:varchar(100);not null"`
	StartTime string `gorm:"type:varchar(100)"`
	// 打卡时长（min）
	Dur int `gorm:"type:int"`
}
