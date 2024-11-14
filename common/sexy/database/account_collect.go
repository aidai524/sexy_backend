package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// AccountCollect 用户收藏表
type AccountCollect struct {
	Model
	Address   string `gorm:"type:varchar(64);not null;column:address;index:idx_address_time,priority:1;uniqueIndex:idx_address_id,priority:1" json:"address"` // 用户地址
	ProjectID uint64 `gorm:"type:bigint;not null;uniqueIndex:idx_address_id,priority:2" json:"project_id"`                                                    // 项目ID
	Time      int64  `gorm:"type:bigint;not null;index:idx_address_time,priority:2" json:"time"`                                                              // 收藏时间
}

func (AccountCollect) TableName() string {
	return "account_collect"
}

// GetAccountCollect 查询
func GetAccountCollect(db *gorm.DB, address string, projectID uint64) (*AccountCollect, error) {
	var project AccountCollect
	err := db.Where("address = ? and project_id = ?", address, projectID).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateAccountCollect 创建
func CreateAccountCollect(db *gorm.DB, projectLike *AccountCollect) error {
	// 使用 GORM 创建记录
	if err := db.Create(projectLike).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// DeleteAccountCollect 根据地址和项目ID删除用户收藏记录
func DeleteAccountCollect(db *gorm.DB, address string, projectID uint64) error {
	if err := db.Where("address = ? AND project_id = ?", address, projectID).Unscoped().Delete(&AccountCollect{}).Error; err != nil {
		return fmt.Errorf("failed to delete account collect: %w", err)
	}
	return nil
}
