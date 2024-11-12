package database

// AccountCollect 用户收藏表
type AccountCollect struct {
	Model
	Address   string `gorm:"type:varchar(64);not null;column:address;index:idx_address_time,priority:1" json:"address"` // 用户地址
	ProjectID uint64 `gorm:"type:bigint;not null" json:"project_id"`                                                    // 项目ID
	Time      int64  `gorm:"type:bigint;not null;index:idx_address_time,priority:2" json:"time"`                        // 收藏时间
}

func (AccountCollect) TableName() string {
	return "account_collect"
}
