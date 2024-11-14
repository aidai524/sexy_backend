package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ProjectSuperLike struct {
	Model
	ProjectID uint64 `gorm:"type:bigint;index;column:project_id" json:"project_id"`
	Address   string `gorm:"type:varchar(64);not null;column:address;index:idx_address_time,priority:1" json:"address"` // 用户地址
	Time      int64  `gorm:"type:bigint;not null;column:time;index:idx_address_time,priority:2" json:"time"`
}

func (ProjectSuperLike) TableName() string {
	return "project_super_like"
}

// GetProjectSuperLike 查询
func GetProjectSuperLike(db *gorm.DB, address string, projectID uint64) (*ProjectSuperLike, error) {
	var project ProjectSuperLike
	err := db.Where("address = ? and project_id = ?", address, projectID).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProjectSuperLike 创建超级点赞
func CreateProjectSuperLike(db *gorm.DB, projectLike *ProjectSuperLike) error {
	// 使用 GORM 创建记录
	if err := db.Create(projectLike).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// GetTodayProjectSuperLikes 查询某个 address 当天的超级喜欢数据
func GetTodayProjectSuperLikes(db *gorm.DB, address string) ([]*ProjectSuperLike, error) {
	var superLikes []*ProjectSuperLike

	// 获取当天的开始时间和结束时间（以 UTC 时间为例）
	t := time.Now()
	endT := t.AddDate(0, 0, 1)
	startT := t.AddDate(0, 0, -1)
	startOfDay := time.Date(startT.Year(), startT.Month(), startT.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
	endOfDay := time.Date(endT.Year(), endT.Month(), endT.Day(), 0, 0, 0, 0, time.Local).UnixMilli()

	// 查询指定 address 当天的记录
	if err := db.Where("address = ? AND time >= ? AND time < ?", address, startOfDay, endOfDay).
		Find(&superLikes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get today's super likes for address %s: %w", address, err)
	}

	return superLikes, nil
}
