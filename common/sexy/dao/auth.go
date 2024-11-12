package dao

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	redis2 "sexy_backend/common/redis"
)

const maxTokensPerUser = 1 // 每个用户最多允许的 Token 数量

// SetAuth 设置 Token，同时确保用户的 Token 数量不超过 N
func SetAuth(pool *redis.Pool, token string, account string) (err error) {
	var (
		conn = pool.Get()
	)
	defer func() {
		_ = conn.Close()
	}()

	// Redis 列表 key，用于存储该用户的所有 Token
	userTokensKey := fmt.Sprintf(redis2.UserTokens, account)

	// 检查当前用户的 Token 数量
	tokenCount, err := redis.Int(conn.Do("LLEN", userTokensKey))
	if err != nil {
		return err
	}

	// 如果超过最大 Token 数量，删除最早的 Token
	if tokenCount >= maxTokensPerUser {
		_, err = conn.Do("LPOP", userTokensKey) // 删除列表中最早的 Token
		if err != nil {
			return err
		}
	}

	// 将新的 Token 添加到用户的 Token 列表中
	_, err = conn.Do("RPUSH", userTokensKey, token)
	if err != nil {
		return err
	}

	// 设置 Token 和账户映射
	_, err = conn.Do("SET", fmt.Sprintf(redis2.AuthToken, token), account)
	return err
}

// GetAuth 获取 Token 对应的账户
func GetAuth(pool *redis.Pool, token string) (account string, err error) {
	var (
		conn = pool.Get()
	)
	defer func() {
		_ = conn.Close()
	}()

	// 获取 Token 对应的账户
	account, err = redis.String(conn.Do("GET", fmt.Sprintf(redis2.AuthToken, token)))
	if errors.Is(err, redis.ErrNil) {
		err = nil
		return
	}
	return
}
