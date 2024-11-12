package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
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
	Like              uint64 `gorm:"type:bigint;not null;default:0;column:like" json:"like"`                                    // 喜欢
	UnLike            uint64 `gorm:"type:bigint;not null;default:0;column:un_like" json:"un_like"`                              // 不喜欢
	SuperLike         uint64 `gorm:"type:bigint;not null;default:0;column:super_like" json:"super_like"`                        // 超级喜欢
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
func GetProjectByTokenName(db *gorm.DB, tokenName string) (*Project, error) {
	var project Project
	err := db.Where("token_name = ?", tokenName).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetProjectByTokenSymbol 根据 TokenSymbol 查询项目
func GetProjectByTokenSymbol(db *gorm.DB, tokenSymbol string) (*Project, error) {
	var project Project
	err := db.Where("token_symbol = ?", tokenSymbol).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find project by TokenSymbol %s: %w", tokenSymbol, err)
	}
	return &project, nil
}

// GetProjectByID 根据 id 查询项目
func GetProjectByID(db *gorm.DB, id uint64) (*Project, error) {
	var project Project
	err := db.Where("id = ?", id).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find project by id %v: %w", id, err)
	}
	return &project, nil
}

// SearchProjects 模糊查询
func SearchProjects(db *gorm.DB, query string, limit, offset int) ([]*Project, error) {
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

	return projects, nil
}

func GetProjectsForAccountOptimized(db *gorm.DB, account string, limit int) ([]*Project, error) {
	var projects []*Project
	boostLimit := limit // 控制 boost 优先部分的查询数量

	// If the account is empty, skip the filtering by like/unlike/superlike
	if account == "" {
		// Fetch all projects, boost priority
		queryBoost := `
            SELECT p.*
            FROM project p
            WHERE p.boost = true
            ORDER BY p.id OFFSET floor(random() * (SELECT count(*) FROM project WHERE boost = true)) LIMIT ?
        `
		if err := db.Raw(queryBoost, boostLimit).Scan(&projects).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch boosted projects: %w", err)
		}

		// If boost projects are not enough, fetch non-boost projects
		if len(projects) < limit {
			remaining := limit - len(projects)
			queryNonBoost := `
                SELECT p.*
                FROM project p
                WHERE p.boost = false
                ORDER BY p.id OFFSET floor(random() * (SELECT count(*) FROM project WHERE boost = false)) LIMIT ?
            `
			var nonBoostProjects []*Project
			if err := db.Raw(queryNonBoost, remaining).Scan(&nonBoostProjects).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch non-boosted projects: %w", err)
			}
			projects = append(projects, nonBoostProjects...)
		}
	} else {
		// Step 1: Fetch boosted projects with account filtering
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
            ORDER BY p.id OFFSET floor(random() * (SELECT count(*) FROM project WHERE boost = true)) LIMIT ?
        `
		if err := db.Raw(queryBoost, account, boostLimit).Scan(&projects).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch boosted projects: %w", err)
		}

		// Step 2: If boost projects are not enough, fetch non-boost projects with account filtering
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
                ORDER BY p.id OFFSET floor(random() * (SELECT count(*) FROM project WHERE boost = false)) LIMIT ?
            `
			var nonBoostProjects []*Project
			if err := db.Raw(queryNonBoost, account, remaining).Scan(&nonBoostProjects).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch non-boosted projects: %w", err)
			}
			projects = append(projects, nonBoostProjects...)
		}
	}

	// Randomize the order of the final results
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(projects), func(i, j int) { projects[i], projects[j] = projects[j], projects[i] })

	return projects, nil
}

// IncrementProjectsLike increments the like counter for the project with the given ID.
func IncrementProjectsLike(db *gorm.DB, projectID uint64) error {
	// Update the like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("like", gorm.Expr("like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment like count: %w", err)
	}
	return nil
}

// IncrementProjectsUnLike increments the un_like counter for the project with the given ID.
func IncrementProjectsUnLike(db *gorm.DB, projectID uint64) error {
	// Update the un_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("un_like", gorm.Expr("un_like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment unlike count: %w", err)
	}
	return nil
}

// IncrementProjectsSuperLike increments the super_like counter for the project with the given ID.
func IncrementProjectsSuperLike(db *gorm.DB, projectID uint64) error {
	// Update the super_like counter for the project by incrementing it by 1.
	if err := db.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("super_like", gorm.Expr("super_like + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment super like count: %w", err)
	}
	return nil
}
