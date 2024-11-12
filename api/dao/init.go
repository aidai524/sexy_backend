package dao

import (
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"sexy_backend/api/conf"
	m "sexy_backend/common/pgsql"
	r "sexy_backend/common/redis"
	"sexy_backend/common/sexy/database"
)

type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:    m.NewGormPostgres(c.Pgsql.DB, database.GetAllDBTable(), nil),
		redis: r.NewRedisPool(c.Redis),
	}
	return
}

func (d *Dao) GetDB() (db *gorm.DB) {
	return d.db
}

func (d *Dao) GetRedis() (redis *redis.Pool) {
	return d.redis
}

func (d *Dao) WithTrx(f func(tx *gorm.DB) (err error)) (err error) {
	return d.db.Transaction(func(tx *gorm.DB) error {
		//// 设置事务隔离级别
		//if err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ").Error; err != nil {
		//	return err
		//}
		return f(tx)
	})
}
