package database

// Account 用户表
type Account struct {
	Model
	Address string `gorm:"type:varchar(64);not null;column:address;uniqueIndex:idx_address" json:"address"` // 用户地址
	VipType string `gorm:"type:varchar(32);not null;column:vip_type" json:"vip_type"`                       // vip等级
	Time    int64  `gorm:"type:bigint;not null;default:0" json:"time"`                                      // 用户注册时间
}

func (Account) TableName() string {
	return "account"
}
