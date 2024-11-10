package database

import "time"

// Account 用户表
type Account struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"-"`
	Address    string    `gorm:"type:varchar(64);not null;column:address;unique" json:"address"` // 用户地址
	VipType    string    `gorm:"type:varchar(32);not null;column:vip_type" json:"vip_type"`      // vip等级
	Time       int64     `gorm:"type:bigint;not null;default:0" json:"time"`                     // 用户注册时间
	CreateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"-"`
	UpdateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

func (Account) TableName() string {
	return "account"
}
