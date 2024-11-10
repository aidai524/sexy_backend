package database

import "time"

type ProjectUnLike struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"-"`
	ProjectID  uint64    `gorm:"type:bigint;index;column:project_id" json:"project_id"`
	CreateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"-"`
	UpdateTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

func (ProjectUnLike) TableName() string {
	return "project_un_like"
}
