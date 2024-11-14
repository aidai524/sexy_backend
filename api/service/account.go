package service

import (
	"sexy_backend/api/model"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/database"
	"time"
)

func (s *Service) GetAccountData(account string) (data *model.AccountData, err error) {
	var (
		timeList         []int64
		accountActive    []*database.AccountActive
		projectSuperLike []*database.ProjectSuperLike
	)

	data = &model.AccountData{
		SuperLikeTotal: 1,
	}

	projectSuperLike, err = database.GetTodayProjectSuperLikes(s.Dao.GetDB(), account)
	if err != nil {
		log.Error("GetAccountData - GetTodayProjectSuperLikes err: %v", err)
		return
	}

	data.SuperLikeAvailable = len(projectSuperLike)

	now := time.Now()
	for i := 0; i < 7; i++ {
		t := now.AddDate(0, 0, -i)
		timeList = append(timeList, time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UnixMilli())
	}

	// 查询用户7天内是否活跃
	accountActive, err = database.GetAccountActiveList(s.Dao.GetDB(), account, timeList)
	if err != nil {
		log.Error("GetAccountActiveList error: %v", err)
		return
	}
	if len(accountActive) == len(timeList) {
		data.SuperLikeTotal = 2
	}
	return
}
