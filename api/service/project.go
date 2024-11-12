package service

import (
	"gorm.io/gorm"
	"sexy_backend/common/ecode"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/database"
	"sexy_backend/common/sexyerror"
	"time"
)

func (s *Service) PostProject(account, tokenName, tokenSymbol, icon, video, backgroundStory, futureDevelopment, whitePaper, x, tg, country string) (err error) {
	var project *database.Project
	project, err = database.GetProjectByTokenName(s.Dao.GetDB(), tokenName)
	if err != nil {
		log.Error("PostProject - GetProjectByTokenName err: %v", err)
		return
	}
	if project != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: "token name already exists error"}
		return
	}

	project, err = database.GetProjectByTokenSymbol(s.Dao.GetDB(), tokenSymbol)
	if err != nil {
		log.Error("PostProject - GetProjectByTokenName err: %v", err)
		return
	}
	if project != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: "token symbol already exists error"}
		return
	}

	err = database.CreateProject(s.Dao.GetDB(), &database.Project{
		Account:           account,
		TokenName:         tokenName,
		TokenSymbol:       tokenSymbol,
		Icon:              icon,
		Video:             video,
		BackgroundStory:   backgroundStory,
		FutureDevelopment: futureDevelopment,
		WhitePaper:        whitePaper,
		X:                 x,
		Tg:                tg,
		Country:           country,
		Time:              time.Now().UnixMilli(),
	})
	if err != nil {
		log.Error("PostProject - CreateProject err: %v", err)
		return
	}
	return
}

func (s *Service) GetProject(id uint64, tokenName, tokenSymbol string) (project *database.Project, err error) {
	if id > 0 {
		project, err = database.GetProjectByID(s.Dao.GetDB(), id)
		if err != nil {
			log.Error("GetProject - GetProjectByID err: %v", err)
			return
		}
		return
	}

	if tokenName != "" {
		project, err = database.GetProjectByTokenName(s.Dao.GetDB(), tokenName)
		if err != nil {
			log.Error("GetProject - GetProjectByTokenName err: %v", err)
			return
		}
		return
	}

	if tokenSymbol != "" {
		project, err = database.GetProjectByTokenSymbol(s.Dao.GetDB(), tokenSymbol)
		if err != nil {
			log.Error("GetProject - GetProjectByTokenSymbol err: %v", err)
			return
		}
		return
	}
	return
}

func (s *Service) GetProjectSearch(text string, limit, offset int) (projects []*database.Project, haxNextPage bool, err error) {
	projects, err = database.SearchProjects(s.Dao.GetDB(), text, limit+1, offset)
	if err != nil {
		log.Error("GetProjectSearch - SearchProjects err: %v", err)
		return
	}

	haxNextPage = len(projects) > limit
	if len(projects) > limit {
		projects = projects[0:limit]
	}
	return
}

func (s *Service) GetProjectList(account string, limit int) (projects []*database.Project, haxNextPage bool, err error) {
	projects, err = database.GetProjectsForAccountOptimized(s.Dao.GetDB(), account, limit+1)
	if err != nil {
		log.Error("GetProjectList - GetProjectsForAccountOptimized err: %v", err)
		return
	}

	haxNextPage = len(projects) > limit
	if len(projects) > limit {
		projects = projects[0:limit]
	}
	return
}

func (s *Service) PostProjectLike(account string, id uint64) (err error) {
	var projectLike *database.ProjectLike
	projectLike, err = database.GetProjectLike(s.Dao.GetDB(), account, id)
	if err != nil {
		log.Error("PostProjectLike - GetProjectLike err: %v", err)
		return
	}

	if projectLike != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: "repeat like error"}
		return
	}

	err = s.Dao.WithTrx(func(tx *gorm.DB) (err error) {
		err = database.CreateProjectLike(tx, &database.ProjectLike{
			Address:   account,
			ProjectID: id,
		})
		if err != nil {
			log.Error("PostProjectLike - CreateProjectLike err: %v", err)
			return
		}

		err = database.IncrementProjectsLike(tx, id)
		if err != nil {
			log.Error("PostProjectLike - IncrementProjectsLike err: %v", err)
			return
		}
		return
	})
	if err != nil {
		log.Error("PostProjectLike - CreateProjectLike err: %v", err)
		return
	}
	return
}

func (s *Service) PostProjectUnLike(account string, id uint64) (err error) {
	var projectUnLike *database.ProjectUnLike
	projectUnLike, err = database.GetProjectUnLike(s.Dao.GetDB(), account, id)
	if err != nil {
		log.Error("PostProjectUnLike - GetProjectUnLike err: %v", err)
		return
	}

	if projectUnLike != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: "repeat un like error"}
		return
	}

	err = s.Dao.WithTrx(func(tx *gorm.DB) (err error) {
		err = database.CreateProjectUnLike(tx, &database.ProjectUnLike{
			Address:   account,
			ProjectID: id,
		})
		if err != nil {
			log.Error("PostProjectUnLike - CreateProjectLike err: %v", err)
			return
		}

		err = database.IncrementProjectsUnLike(tx, id)
		if err != nil {
			log.Error("PostProjectUnLike - IncrementProjectsLike err: %v", err)
			return
		}
		return
	})
	if err != nil {
		log.Error("PostProjectUnLike - CreateProjectLike err: %v", err)
		return
	}
	return
}

func (s *Service) PostProjectSuperLike(account string, id uint64) (err error) {
	var projectSuperLike *database.ProjectSuperLike
	projectSuperLike, err = database.GetProjectSuperLike(s.Dao.GetDB(), account, id)
	if err != nil {
		log.Error("PostProjectSuperLike - GetProjectSuperLike err: %v", err)
		return
	}

	if projectSuperLike != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: "repeat super like error"}
		return
	}

	err = s.Dao.WithTrx(func(tx *gorm.DB) (err error) {
		err = database.CreateProjectSuperLike(tx, &database.ProjectSuperLike{
			Address:   account,
			ProjectID: id,
		})
		if err != nil {
			log.Error("PostProjectSuperLike - CreateProjectUnLike err: %v", err)
			return
		}

		err = database.IncrementProjectsSuperLike(tx, id)
		if err != nil {
			log.Error("PostProjectSuperLike - IncrementProjectsUnLike err: %v", err)
			return
		}
		return
	})
	if err != nil {
		log.Error("PostProjectSuperLike - IncrementProjectsUnLike err: %v", err)
		return
	}
	return
}
