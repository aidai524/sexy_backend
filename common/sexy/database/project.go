package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"sexy_backend/common/log"
	"time"
)

type Project struct {
	Model
	Account           string `gorm:"type:varchar(64);not null;column:account" json:"account"`
	TokenName         string `gorm:"type:varchar(100);not null;uniqueIndex:uidx_name;column:token_name" json:"token_name"`      // 项目代币名称(唯一)
	TokenSymbol       string `gorm:"type:varchar(20);not null;uniqueIndex:uidx_symbol;column:token_symbol" json:"token_symbol"` // 项目代币symbol(唯一)
	Icon              string `gorm:"type:varchar(300);column:icon" json:"icon"`                                                 // 项目图标URL
	Video             string `gorm:"type:varchar(300);column:video" json:"video"`                                               // 视频URL
	BackgroundStory   string `gorm:"type:text;column:background_story" json:"background_story"`                                 // 背景故事
	FutureDevelopment string `gorm:"type:text;column:future_development" json:"future_development"`                             // 未来发展
	WhitePaper        string `gorm:"type:varchar(300);column:white_paper" json:"white_paper"`                                   // 白皮书URL
	X                 string `gorm:"type:varchar(300);column:x" json:"x"`                                                       // 项目方的社交媒体信息
	Tg                string `gorm:"type:varchar(300);column:tg" json:"tg"`                                                     // 项目方的社交媒体信息
	Country           string `gorm:"type:varchar(100);column:country" json:"country"`                                           // 国家信息
	Like              uint64 `gorm:"type:bigint;not null;default:0;column:like" json:"like"`                                    // 喜欢数量
	UnLike            uint64 `gorm:"type:bigint;not null;default:0;column:un_like" json:"un_like"`                              // 不喜欢数量
	SuperLike         uint64 `gorm:"type:bigint;not null;default:0;column:super_like" json:"super_like"`                        // 超级喜欢数量
	Collect           uint64 `gorm:"type:bigint;not null;default:0;column:collect" json:"collect"`                              // 收藏数量
	IsCollect         bool   `gorm:"-" json:"is_collect"`                                                                       // 是否收藏
	Boost             bool   `gorm:"type:boolean;not null;default:0" json:"boost"`                                              // 是否提供曝光度
	BoostTime         int64  `gorm:"type:bigint;not null;default:0" json:"boost_time"`                                          // 是否提供曝光度
	Time              int64  `gorm:"type:bigint;not null;default:0" json:"time"`                                                // 项目创建时间
}

func (Project) TableName() string {
	return "project"
}

// CreateProject 创建项目
func CreateProject(db *gorm.DB, project *Project) error {
	// 使用 GORM 创建记录
	if err := db.Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// GetProjectByTokenName 根据 TokenName 查询项目
func GetProjectByTokenName(db *gorm.DB, account string, tokenName string) (*Project, error) {
	var project Project
	err := db.Where("token_name = ?", tokenName).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = GetProjectAccountIsCollect(db, account, []*Project{&project})
	if err != nil {
		log.Error("GetProjectByTokenName - err: %v", err)
		return nil, err
	}
	return &project, nil
}

// GetProjectByTokenSymbol 根据 TokenSymbol 查询项目
func GetProjectByTokenSymbol(db *gorm.DB, account, tokenSymbol string) (*Project, error) {
	var project Project
	err := db.Where("token_symbol = ?", tokenSymbol).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find project by TokenSymbol %s: %w", tokenSymbol, err)
	}
	err = GetProjectAccountIsCollect(db, account, []*Project{&project})
	if err != nil {
		log.Error("GetProjectByTokenName - err: %v", err)
		return nil, err
	}
	return &project, nil
}

// GetProjectByID 根据 id 查询项目
func GetProjectByID(db *gorm.DB, account string, id uint64) (*Project, error) {
	var project Project
	err := db.Where("id = ?", id).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find project by id %v: %w", id, err)
	}
	err = GetProjectAccountIsCollect(db, account, []*Project{&project})
	if err != nil {
		log.Error("GetProjectByTokenName - err: %v", err)
		return nil, err
	}
	return &project, nil
}

// SearchProjects 模糊查询
func SearchProjects(db *gorm.DB, account, query string, limit, offset int) ([]*Project, error) {
	var projects []*Project

	// 模糊查询条件，分页设置，以及时间倒序排序
	err := db.Where("token_name LIKE ? OR token_symbol LIKE ?", query+"%", query+"%").
		Order("time DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}

	err = GetProjectAccountIsCollect(db, account, projects)
	if err != nil {
		log.Error("GetProjectByTokenName - err: %v", err)
		return nil, err
	}
	return projects, nil
}
func GetProjectsForAccountOptimized(db *gorm.DB, account string, limit int) ([]*Project, error) {
	var projects []*Project
	boostLimit := limit

	if account == "" {
		// Boosted projects without account filtering
		queryBoost := `
            SELECT p.*
            FROM project p
            WHERE p.boost = true
            ORDER BY p.id LIMIT ?
        `
		if err := db.Raw(queryBoost, boostLimit*2).Scan(&projects).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch boosted projects: %w", err)
		}

		// If not enough boosted projects, fetch non-boosted ones
		if len(projects) < limit {
			remaining := limit - len(projects)
			queryNonBoost := `
                SELECT p.*
                FROM project p
                WHERE p.boost = false
                ORDER BY p.id LIMIT ?
            `
			var nonBoostProjects []*Project
			if err := db.Raw(queryNonBoost, remaining*2).Scan(&nonBoostProjects).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch non-boosted projects: %w", err)
			}
			projects = append(projects, nonBoostProjects...)
		}
	} else {
		// Boosted projects with account filtering
		queryBoost := `
            SELECT p.*
            FROM project p
            LEFT JOIN project_like pl ON p.id = pl.project_id AND pl.address = ?
            LEFT JOIN project_un_like pu ON p.id = pu.project_id AND pu.address = ?
            LEFT JOIN project_super_like ps ON p.id = ps.project_id AND ps.address = ?
            WHERE pl.project_id IS NULL
              AND pu.project_id IS NULL
              AND ps.project_id IS NULL
              AND p.boost = true
            ORDER BY p.id LIMIT ?
        `
		if err := db.Raw(queryBoost, account, account, account, boostLimit*2).Scan(&projects).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch boosted projects: %w", err)
		}

		// Non-boosted projects with account filtering
		if len(projects) < limit {
			remaining := limit - len(projects)
			queryNonBoost := `
                SELECT p.*
                FROM project p
                LEFT JOIN project_like pl ON p.id = pl.project_id AND pl.address = ?
                LEFT JOIN project_un_like pu ON p.id = pu.project_id AND pu.address = ?
                LEFT JOIN project_super_like ps ON p.id = ps.project_id AND ps.address = ?
                WHERE pl.project_id IS NULL
                  AND pu.project_id IS NULL
                  AND ps.project_id IS NULL
                  AND p.boost = false
                ORDER BY p.id LIMIT ?
            `
			var nonBoostProjects []*Project
			if err := db.Raw(queryNonBoost, account, account, account, remaining*2).Scan(&nonBoostProjects).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch non-boosted projects: %w", err)
			}
			projects = append(projects, nonBoostProjects...)
		}
	}

	// Randomize order and limit to the required count
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(projects), func(i, j int) { projects[i], projects[j] = projects[j], projects[i] })

	if len(projects) > limit {
		projects = projects[:limit]
	}

	err := GetProjectAccountIsCollect(db, account, projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// IncrementProjectsLike  like
func IncrementProjectsLike(db *gorm.DB, projectID uint64) error {
	// Wrap the `like` field in double quotes to avoid syntax issues with the SQL keyword
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("\"like\"", gorm.Expr("\"like\" + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment like count: %w", err)
	}
	return nil
}

// IncrementProjectsUnLike  添加UnLike数量
func IncrementProjectsUnLike(db *gorm.DB, projectID uint64) error {
	// Update the un_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("un_like", gorm.Expr("un_like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment unlike count: %w", err)
	}
	return nil
}

// IncrementProjectsSuperLike 添加superLike数量
func IncrementProjectsSuperLike(db *gorm.DB, projectID uint64) error {
	// Update the super_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("super_like", gorm.Expr("super_like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment super like count: %w", err)
	}
	return nil
}

// IncrementProjectsCollect 添加收藏数量
func IncrementProjectsCollect(db *gorm.DB, projectID uint64) error {
	// Update the super_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("super_like", gorm.Expr("super_like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment super like count: %w", err)
	}
	return nil
}

// ReductionProjectsCollect 减少收藏数量
func ReductionProjectsCollect(db *gorm.DB, projectID uint64) error {
	// Update the super_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("super_like", gorm.Expr("super_like - ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment super like count: %w", err)
	}
	return nil
}

// GetAccountCollectList 查询用户的收藏列表
func GetAccountCollectList(db *gorm.DB, address string, limit int, offset int) ([]*Project, error) {
	var projects []*Project

	// 查询用户的收藏列表，并联接到项目表
	query := db.Table("account_collect as ac").
		Select("p.*, true as is_collect").
		Joins("JOIN project p ON ac.project_id = p.id").
		Where("ac.address = ?", address).
		Order("ac.time DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Scan(&projects).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user collect list: %w", err)
	}

	for _, project := range projects {
		project.IsCollect = true
	}
	return projects, nil
}

func GetProjectAccountIsCollect(db *gorm.DB, account string, projects []*Project) error {
	// 如果传入的 projects 为空，直接返回
	if len(projects) == 0 {
		return nil
	}

	if len(account) == 0 {
		return nil
	}

	// 获取用户收藏的项目 ID 列表
	var collectProjectIDs []uint64
	if err := db.Model(&AccountCollect{}).
		Where("address = ?", account).
		Pluck("project_id", &collectProjectIDs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to fetch collected projects for account %s: %w", account, err)
	}

	// 将收藏的项目 ID 存入一个 map，以便于快速查找
	collectMap := make(map[uint64]struct{}, len(collectProjectIDs))
	for _, projectID := range collectProjectIDs {
		collectMap[projectID] = struct{}{}
	}

	// 遍历项目列表，判断每个项目是否在收藏的项目列表中
	for _, project := range projects {
		if _, exists := collectMap[project.ID]; exists {
			project.IsCollect = true
		} else {
			project.IsCollect = false
		}
	}

	return nil
}
