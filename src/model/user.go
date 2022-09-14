package model

import "time"

type User struct {
	UserId     string    `json:"UserId" gorm:"column:user_id;type:bigint(20) unsigned not null AUTO_INCREMENT;primaryKey;"`
	LoginName  string    `json:"LoginName" gorm:"column:login_name"`
	UserName   string    `json:"UserName" gorm:"column:user_name"`
	Email      string    `json:"Email" gorm:"column:email"`
	UpdateTime time.Time `json:"UpdateTime" gorm:"column:update_time"`
}

func (User) TableName() string {
	return "sys_user"
}
