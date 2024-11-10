package dao

import (
	"gorm.io/gorm"
	"sexy_backend/api/conf"
	m "sexy_backend/common/pgsql"
	"sexy_backend/common/sexy/database"
)

type Dao struct {
	db *gorm.DB
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: m.NewGormPostgres(c.Pgsql.DB, database.GetAllDBTable(), nil),
	}
	return
}
