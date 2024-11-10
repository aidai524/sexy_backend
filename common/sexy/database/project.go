package database

import "time"

type Project struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"-"`
	Account           string    `gorm:"type:varchar(64);not null;column:account" json:"account"`
	TokenName         string    `gorm:"type:varchar(100);not null;unique;column:token_name" json:"token_name"`    // 项目代币名称(唯一)
	TokenSymbol       string    `gorm:"type:varchar(20);not null;unique;column:token_symbol" json:"token_symbol"` // 项目代币symbol(唯一)
	Icon              string    `gorm:"type:varchar(300);column:icon" json:"icon"`                                // 项目图标URL
	Video             string    `gorm:"type:varchar(300);column:video" json:"video"`                              // 视频URL
	BackgroundStory   string    `gorm:"type:text;column:background_story" json:"background_story"`                // 背景故事
	FutureDevelopment string    `gorm:"type:text;column:future_development" json:"future_development"`            // 未来发展
	WhitePaper        string    `gorm:"type:varchar(300);column:white_paper" json:"white_paper"`                  // 白皮书URL
	X                 string    `gorm:"type:varchar(300);column:x" json:"x"`                                      // 项目方的社交媒体信息
	Tg                string    `gorm:"type:varchar(300);column:tg" json:"tg"`                                    // 项目方的社交媒体信息
	Country           string    `gorm:"type:varchar(100);column:country" json:"country"`                          // 国家信息
	Like              uint64    `gorm:"type:bigint;not null;default:0;column:like" json:"like"`                   // 喜欢
	UnLike            uint64    `gorm:"type:bigint;not null;default:0;column:un_like" json:"un_like"`             // 不喜欢
	SuperLike         uint64    `gorm:"type:bigint;not null;default:0;column:super_like" json:"super_like"`       // 超级喜欢
	Boost             bool      `gorm:"type:tinyint(1);not null;default:0" json:"boost"`                          // 是否提供曝光度
	BoostTime         int64     `gorm:"type:bigint;not null;default:0" json:"boost_time"`                         // 是否提供曝光度
	Time              int64     `gorm:"type:bigint;not null;default:0" json:"time"`                               // 项目创建时间
	CreateTime        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"-"`
	UpdateTime        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

func (Project) TableName() string {
	return "project"
}
