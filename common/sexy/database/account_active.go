package database

import "time"

// AccountActive 用户活跃表
type AccountActive struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"-"`
	Address    string    `gorm:"type:varchar(64);not null;column:address;uniqueIndex:uidx_address_time,priority:1" json:"address"`
	Time       int64     `gorm:"type:bigint;not null;uniqueIndex:uidx_address_time,priority:2" json:"time"`
	CreateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"-"`
	UpdateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

func (AccountActive) TableName() string {
	return "account_active"
}
