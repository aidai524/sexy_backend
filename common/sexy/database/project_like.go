package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProjectLike struct {
	Model
	ProjectID uint64 `gorm:"type:bigint;index;column:project_id" json:"project_id"`
	Address   string `gorm:"type:varchar(64);not null;column:address" json:"address"` // 用户地址
}

func (ProjectLike) TableName() string {
	return "project_like"
}

// GetProjectLike 查询
func GetProjectLike(db *gorm.DB, address string, projectID uint64) (*ProjectLike, error) {
	var project ProjectLike
	err := db.Where("account = ? and project_id = ?", address, projectID).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProjectLike 创建like
func CreateProjectLike(db *gorm.DB, projectLike *ProjectLike) error {
	// 使用 GORM 创建记录
	if err := db.Create(projectLike).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}
