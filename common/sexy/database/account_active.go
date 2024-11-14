package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AccountActive 用户活跃表
type AccountActive struct {
	Model
	Address string `gorm:"type:varchar(64);not null;column:address;uniqueIndex:uidx_address_time,priority:1" json:"address"`
	Time    int64  `gorm:"type:bigint;not null;uniqueIndex:uidx_address_time,priority:2" json:"time"`
}

func (AccountActive) TableName() string {
	return "account_active"
}

// CreateAccountActive 插入用户活跃记录
func CreateAccountActive(db *gorm.DB, address string, time int64) error {
	accountActive := AccountActive{
		Address: address,
		Time:    time,
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}, {Name: "time"}}, // 指定发生冲突的列
		DoNothing: true,
	}).Create(&accountActive).Error
}

// GetAccountActiveList 根据 address 和 time 列表查询用户活跃记录
func GetAccountActiveList(db *gorm.DB, address string, times []int64) ([]*AccountActive, error) {
	var accountActives []*AccountActive

	// 查询符合 address 列表和 time 列表条件的数据
	if err := db.Where("address = ? AND time IN ?", address, times).Find(&accountActives).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account active records: %w", err)
	}

	return accountActives, nil
}
