package database

// AccountActive 用户活跃表
type AccountActive struct {
	Model
	Address string `gorm:"type:varchar(64);not null;column:address;uniqueIndex:uidx_address_time,priority:1" json:"address"`
	Time    int64  `gorm:"type:bigint;not null;uniqueIndex:uidx_address_time,priority:2" json:"time"`
}

func (AccountActive) TableName() string {
	return "account_active"
}
