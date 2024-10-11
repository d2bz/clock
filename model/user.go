package model

import "gorm.io/gorm"

type User struct {
	// ToDo 添加UUID，避免用户与Telephone绑定
	gorm.Model
	UserID    string
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Avatar    string
	// Durations []Duration `gorm:"foreignKey:Tel;references:Telephone"`
}

type SimpleUser struct {
	Username  string
	TimeTotal int
}
