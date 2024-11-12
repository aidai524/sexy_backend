package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProjectSuperLike struct {
	Model
	ProjectID uint64 `gorm:"type:bigint;index;column:project_id" json:"project_id"`
	Address   string `gorm:"type:varchar(64);not null;column:address" json:"address"` // 用户地址
}

func (ProjectSuperLike) TableName() string {
	return "project_super_like"
}

// GetProjectSuperLike 查询
func GetProjectSuperLike(db *gorm.DB, address string, projectID uint64) (*ProjectSuperLike, error) {
	var project ProjectSuperLike
	err := db.Where("account = ? and project_id = ?", address, projectID).First(&project).Error
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
