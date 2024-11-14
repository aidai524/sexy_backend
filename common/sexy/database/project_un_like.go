package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProjectUnLike struct {
	Model
	ProjectID uint64 `gorm:"type:bigint;index;column:project_id" json:"project_id"`
	Address   string `gorm:"type:varchar(64);not null;column:address;index:idx_address_time,priority:1" json:"address"` // 用户地址
	Time      int64  `gorm:"type:bigint;not null;column:time;index:idx_address_time,priority:2" json:"time"`
}

func (ProjectUnLike) TableName() string {
	return "project_un_like"
}

// GetProjectUnLike 查询
func GetProjectUnLike(db *gorm.DB, address string, projectID uint64) (*ProjectUnLike, error) {
	var project ProjectUnLike
	err := db.Where("address = ? and project_id = ?", address, projectID).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProjectUnLike 创建
func CreateProjectUnLike(db *gorm.DB, projectLike *ProjectUnLike) error {
	// 使用 GORM 创建记录
	if err := db.Create(projectLike).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}
