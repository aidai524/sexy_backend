package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sexy_backend/common/conf"
	"sexy_backend/common/log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPgsql(c *conf.Pgsql) (db *sql.DB) {
	db, err := open(c)
	if err != nil {
		panic(err)
	}
	return
}

func NewGormPostgres(c *conf.Pgsql, dst []interface{}, setSharding func(db *gorm.DB) (err error)) (db *gorm.DB) {
	// 使用 PostgreSQL 驱动打开连接
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info), // 启用调试模式
	})
	if err != nil {
		panic(err)
	}

	// 配置分片函数（如果需要）
	if setSharding != nil {
		err = setSharding(db)
		if err != nil {
			log.Error("setSharding error: %v", err)
			return
		}
	}

	// 配置连接池设置
	pgDB, err := db.DB()
	if err != nil {
		log.Error("[PANIC][Database - Init] Get pgDB error: %s", err)
		return
	}
	pgDB.SetMaxIdleConns(c.Idle)
	pgDB.SetMaxOpenConns(c.Active)
	pgDB.SetConnMaxLifetime(time.Second * time.Duration(c.IdleTimeout))
	log.Info("[Database - Init] database client setting done")

	// 自动迁移数据库表结构
	if !c.NotAutoMigrate {
		err = db.AutoMigrate(dst...)
		if err != nil {
			panic(err)
		}
	}
	log.Info("[Database - Init] database schema auto migrated successfully")
	return db
}

func open(c *conf.Pgsql) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", c.DSN)
	if err != nil {
		log.Error("error opening pgsql %+v: %v", c, err)
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	err = db.Ping()
	return
}

func WithTrx(db *sql.DB, f func(tx *sql.Tx) (err error)) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Error("error getting transaction: %v", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("error exec trx : %v", r)
			_ = tx.Rollback()
			err = fmt.Errorf("%v", r)
		}
	}()

	err = f(tx)
	if err != nil {
		log.Error("error trx: %v, roll back", err)
		if e := tx.Rollback(); e != nil {
			log.Error("error rolling back: %v", e)
		}
		return
	}
	if err := tx.Commit(); err != nil {
		log.Error("error committing: %v", err)
	}
	return
}
