package redis

import (
	"context"
	"encoding/json"
	"math/rand"
	"sexy_backend/common/conf"
	"sexy_backend/common/log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var (
	NormalLockTimeout   = 5
	WriteLockTimeout_NX = 660
	WriteLockTimeout_XX = 5
	ReadLockTimeout     = 600 * time.Second

	NormalLockPrefix = "lock:"
	ReadLockPrefix   = "readLock:"
	WriteLockPrefix  = "writeLock:"

	GetReadLockScript = `if redis.call("EXISTS", KEYS[1]) == 1 then
		return 0
	else
		redis.call("INCR", KEYS[2])
	return 1
	end
`
)

func NewRedisPool(c *conf.Redis) *redis.Pool {
	return &redis.Pool{
		MaxActive:   c.Active,
		MaxIdle:     c.Idle,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(c.Proto, c.Addr, redis.DialDatabase(c.DB), redis.DialPassword(c.Password))
		},
	}
}

/*
WithNewLockTimeout 锁
wait: 已经被锁,是否等待
return: 是否正常执行 false: 被锁, true: 未被锁
*/
func WithNewLockTimeout(pool *redis.Pool, lock string, wait func() bool, f func()) bool {
	value, isWait := getNewLock(pool, lock, NormalLockTimeout, wait)
	if isWait {
		ctx, cancel := context.WithCancel(context.Background())
		go updateNormalLock(pool, lock, NormalLockTimeout, value, ctx)
		f()
		cancel()
		return true
	}
	return false
}

func WithNewReadLockTimeout(pool *redis.Pool, lock string, wait func() bool, f func()) bool {
	isWait := getReadLock(pool, lock, wait)
	if isWait {
		// log.Info("Test----------------getReadLock")
		f()
		releaseReadLock(pool, lock)
		// log.Info("Test----------------releaseReadLock")
		return true
	}
	return false
}

func WithNewWriteLockTimeout(pool *redis.Pool, lock string, wait func() bool, f func()) bool {
	value, isWait := getWriteLock(pool, lock, WriteLockTimeout_NX, wait)
	if isWait {
		ctx, cancel := context.WithCancel(context.Background())
		go updateWriteLock(pool, lock, WriteLockTimeout_XX, value, ctx)
		f()
		cancel()
		// log.Info("Test----------------releaseWriteLock")
		return true
	}
	return false
}

func updateWriteLock(pool *redis.Pool, lock string, timeoutSeconds int, value string, ctx context.Context) {
	key := WriteLockPrefix + lock
	updateLock(pool, key, timeoutSeconds, value, ctx)
}

func updateNormalLock(pool *redis.Pool, lock string, timeoutSeconds int, value string, ctx context.Context) {
	key := NormalLockPrefix + lock
	updateLock(pool, key, timeoutSeconds, value, ctx)
}

func updateLock(pool *redis.Pool, key string, timeoutSeconds int, value string, ctx context.Context) {
	var (
		ts = time.Now()
	)
	for {
		select {
		case <-ctx.Done():
			releaseLockError(pool, key, value)
			return
		default:
			conn := pool.Get()
			_, err := redis.String(conn.Do("SET", key, value, "EX", timeoutSeconds, "XX"))
			_ = conn.Close()
			if err != nil {
				return
			}
			if time.Since(ts) > time.Second {
				log.Info("updateLock, key: %v", key)
				ts = time.Now()
			}
			time.Sleep(time.Duration(1+rand.Intn(500)) * time.Millisecond)
		}
	}
}

func getReadLock(pool *redis.Pool, lock string, wait func() bool) (isLock bool) {
	var (
		readKey           = ReadLockPrefix + lock
		writeKey          = WriteLockPrefix + lock
		ts                = time.Now()
		err               error
		conn              redis.Conn
		success           int // 1 for true, 0 for false
		getReadLockScript = redis.NewScript(2, GetReadLockScript)
	)

	for {
		conn = pool.Get()
		success, err = redis.Int(getReadLockScript.Do(conn, writeKey, readKey))
		_ = conn.Close()
		if err != nil {
			log.Info("getReadLock - get read lock error, readKey: %v, error: %v", readKey, err)
		}
		if success == 1 {
			isLock = true
			return
		}
		// log.Info("getReadLock - writeLock is occupied, readKey: %v", readKey)
		isLock = wait()
		if isLock {
			if time.Since(ts) > time.Second {
				ts = time.Now()
			}
			time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
		} else {
			return
		}
	}
}

func getWriteLock(pool *redis.Pool, lock string, timeoutSeconds int, wait func() bool) (value string, isLock bool) {
	var (
		readKey   = ReadLockPrefix + lock
		writeKey  = WriteLockPrefix + lock
		reply     string
		readCount int
		ts        = time.Now()
		conn      redis.Conn
		err       error
	)
	value = uuid.New().String()
	// log.Info("getWriteLock - init value, writeKey: %v, writeValue: %v", lock, value)

	for {
		conn = pool.Get()
		reply, err = redis.String(conn.Do("SET", writeKey, value, "EX", timeoutSeconds, "NX"))
		_ = conn.Close()
		if err == nil {
			// log.Info("Test----------------getWriteLock")
			for {
				readCount, err = getRedisLockInt(pool, readKey)
				if err == nil && readCount == 0 {
					isLock = true
					return
				}
				if err != nil {
					log.Info("getWriteLock - get readCount error, readKey: %v, error: %v", readKey, err)
				}
				if time.Since(ts) > ReadLockTimeout {
					for {
						conn = pool.Get()
						_, err = redis.String(conn.Do("SET", readKey, "0"))
						_ = conn.Close()
						if err != nil {
							log.Info("getWriteLock - reset readCount error, readKey: %v, error: %v", readKey, err)
							time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
						} else {
							log.Info("getWriteLock - reset readCount to 0, readKey: %v", readKey)
							break
						}
					}
				}
				time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
			}
		} else {
			isLock = wait()
			if isLock {
				if time.Since(ts) > time.Second {
					ts = time.Now()
					if err == redis.ErrNil {
						// log.Info("getWriteLock - writeLock is occupied, writeKey: %v", writeKey)
					} else {
						log.Info("getWriteLock - set writeValue failed, writeKey: %v, reply: %v, error: %v", writeKey, reply, err)
					}
				}
				time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
			} else {
				return
			}
		}
	}
}

func releaseReadLock(pool *redis.Pool, lock string) {
	var (
		readKey = ReadLockPrefix + lock
		err     error
	)
	for {
		conn := pool.Get()
		_, err = redis.Int(conn.Do("DECR", readKey))
		_ = conn.Close()
		if err != nil {
			log.Error("Error release read lock %v: %v", readKey, err)
			time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
		} else {
			return
		}
	}
}

func getNewLock(pool *redis.Pool, lock string, timeoutSeconds int, wait func() bool) (value string, isLock bool) {
	var (
		key = NormalLockPrefix + lock
		ts  = time.Now()
		err error
	)
	value = uuid.New().String()
	log.Info("lock redis: %v, uuid: %v", lock, value)
	for {
		conn := pool.Get()
		_, err = redis.String(conn.Do("SET", key, value, "EX", timeoutSeconds, "NX"))
		_ = conn.Close()
		if err == nil {
			isLock = true
			return
		}
		isLock = wait()
		if isLock {
			if time.Since(ts) > time.Second {
				log.Info("getting lock %v: %v", lock, err)
				ts = time.Now()
			}
			time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
		} else {
			return
		}
	}
}

// WithLockTimeout
// Deprecated: 被废弃请使用 - WithLockTimeoutError
func WithLockTimeout(pool *redis.Pool, lock string, timeoutSeconds int, f func()) {
	conn := pool.Get()
	defer func() { _ = conn.Close() }()
	value := getLock(conn, lock, timeoutSeconds)
	f()
	releaseLock(conn, lock, value)
}

// WithLock
// Deprecated: 被废弃请使用 - WithLockError
func WithLock(pool *redis.Pool, lock string, f func()) {
	WithLockTimeout(pool, lock, 5, f)
}

// getLock
// Deprecated: 被废弃请使用 - getLockError
func getLock(conn redis.Conn, lock string, timeoutSeconds int) (value string) {
	var (
		key = NormalLockPrefix + lock
		ts  = time.Now()
		err error
	)
	value = uuid.New().String()

	for {
		_, err = redis.String(conn.Do("SET", key, value, "EX", timeoutSeconds, "NX"))
		if err == nil {
			return
		}
		if time.Since(ts) > time.Second {
			log.Warn("Error getting lock %v: %v", lock, err)
			ts = time.Now()
		}
		time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
	}
}

// releaseLock warning: release is not absolutely save, e.g.
// 1. ThreadA got the value corresponding to key, and equals to expected value
// 2. key is expired
// 3. ThreadB set a new value to key
// 4. ThreadA finishes, and deletes key(the value is different now)
// So lua script is needed for safe releasing
// Deprecated: 被废弃请使用 - releaseLockError
func releaseLock(conn redis.Conn, lock, value string) {
	var (
		key   = NormalLockPrefix + lock
		err   error
		reply string
	)

	for {
		reply, err = redis.String(conn.Do("GET", key))
		if err == nil {
			if reply == value {
				_, _ = conn.Do("DEL", key)
			} else {
				log.Warn("lock %v value changed before release", key)
			}
			return
		} else if err == redis.ErrNil {
			// 被解锁了
			log.Warn("lock %v not exist before release", key)
			return
		} else {
			log.Error("Error releasing lock %v: %v", key, err)
			time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
		}
	}
}

func WithLockError(pool *redis.Pool, lock string, f func()) (err error) {
	err = WithLockTimeoutError(pool, lock, 5, f)
	return
}

func WithLockTimeoutError(pool *redis.Pool, lock string, timeoutSeconds int, f func()) (err error) {
	value, err := getLockError(pool, lock, timeoutSeconds)
	if err != nil {
		return
	}
	f()
	releaseNormalLockError(pool, lock, value)
	return
}

func getLockError(pool *redis.Pool, lock string, timeoutSeconds int) (value string, err error) {
	var (
		key = NormalLockPrefix + lock
		ts  = time.Now()
	)
	value = uuid.New().String()
	for {
		conn := pool.Get()
		_, err = redis.String(conn.Do("SET", key, value, "EX", timeoutSeconds, "NX"))
		_ = conn.Close()
		if err == nil {
			return
		} else if err == redis.ErrPoolExhausted {
			// 连接池被耗尽
			return
		}
		if time.Since(ts) > time.Second {
			log.Warn("Error getting lock %v: %v", lock, err)
			ts = time.Now()
		}
		time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
	}
}

// func releaseWriteLock(pool *redis.Pool, lock, value string) {
// 	var WriteKey = "writeLock:" + lock
// 	releaseLockError(pool, WriteKey, value)
// }

func releaseNormalLockError(pool *redis.Pool, lock, value string) {
	key := NormalLockPrefix + lock
	releaseLockError(pool, key, value)
}

func releaseLockError(pool *redis.Pool, key, value string) {
	var (
		err   error
		reply string
	)
	for {
		reply, err = getRedisLock(pool, key)
		if err == nil {
			if reply == value {
				delRedisLock(pool, key)
			} else {
				log.Warn("lock %v value changed before release", key)
			}
			return
		} else if err == redis.ErrNil {
			// 被解锁了
			log.Warn("lock %v not exist before release", key)
			return
		} else {
			log.Error("Error releasing lock %v: %v", key, err)
			time.Sleep(time.Duration(1+rand.Intn(100)) * time.Millisecond)
		}
	}
}

func getRedisLock(pool *redis.Pool, key string) (string, error) {
	var (
		conn = pool.Get()
	)
	defer func() {
		_ = conn.Close()
	}()
	return redis.String(conn.Do("GET", key))
}

func getRedisLockInt(pool *redis.Pool, key string) (int, error) {
	var (
		conn = pool.Get()
	)
	defer func() {
		_ = conn.Close()
	}()
	return redis.Int(conn.Do("GET", key))
}

func delRedisLock(pool *redis.Pool, key string) {
	var (
		conn = pool.Get()
	)
	defer func() {
		_ = conn.Close()
	}()
	_, _ = conn.Do("DEL", key)
}

func SetStruct(pool *redis.Pool, key string, obj interface{}) (err error) {
	var (
		conn = pool.Get()
	)
	defer func() { _ = conn.Close() }()
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	_, err = conn.Do("SET", key, data)
	if err != nil {
		return
	}
	return
}

func GetStruct(pool *redis.Pool, key string, obj interface{}, h func(err error)) {
	var (
		b    []byte
		err  error
		conn = pool.Get()
	)
	b, err = redis.Bytes(conn.Do("GET", key))
	_ = conn.Close()
	if err != nil {
		h(err)
		return
	}

	err = json.Unmarshal(b, &obj)
	if err != nil {
		h(err)
		return
	}
	h(nil)
}
