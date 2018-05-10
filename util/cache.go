package util

import (
	"strings"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	po *redis.Pool
	KeyPrefix = ""
)

//func init() {
//	po = NewRedisPool()
//}
type RedisClient struct {
	Conn redis.Conn
	err error
	
	HasPrefix bool
}

func NewRedisPool() (*redis.Pool) {
	addr := strings.Join([]string{
		GetConfig("redis", "address"),
		GetConfig("redis","port")}, ":")
	rp := &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			db, _ := strconv.Atoi(GetConfig("redis", "db"))
			c, err := redis.Dial(
				"tcp",
				addr,
				redis.DialDatabase(db),
				redis.DialPassword(GetConfig("redis","password")),
				redis.DialReadTimeout(10*time.Second),
				redis.DialWriteTimeout(10*time.Second),
				)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}

	return rp
}

func GetRedis() *RedisClient {
	if po == nil || ping() != nil {
		po = NewRedisPool()
	}
	
	return &RedisClient{Conn:po.Get(), err:nil, HasPrefix:false}
}

func (r *RedisClient) WithPrefix(p string) *RedisClient {
	KeyPrefix = p
	r.HasPrefix = true
	return r
}

func ping() error {
	_, err := po.Get().Do("PING")
	return err
}

func (r *RedisClient) SET(key string, val interface{}) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	args := redis.Args{}.Add(key, val)

	_, err := redis.String(r.Conn.Do("SET", args...))
	return err
}

func (r *RedisClient) SETEX(key string, expireSeconds int, val interface{}) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	args := redis.Args{}.Add(key, expireSeconds,val)

	_, err := redis.String(r.Conn.Do("SETEX", args...))
	return err
}

func (r *RedisClient) GET(key string) string {
	if r.err != nil {
		return ""
	}

	key = r.key(key)

	val, err := redis.String(r.Conn.Do("GET", key))
	if err != nil {
		return ""
	}

	return val
}

func (r *RedisClient) EXPIRE(key string, expireSeconds int) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	_, err := redis.Int(r.Conn.Do("EXPIRE", key, expireSeconds))

	return err
}

func (r *RedisClient) DEL(key string) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	_, err := redis.Int(r.Conn.Do("DEL", key))

	return err
}

func (r *RedisClient) HSET(key, field, val string) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	_, err := redis.Int(r.Conn.Do("HSET", key, field, val))
	return err
}

func (r *RedisClient) HGET(key, field string) (string, error) {
	if r.err != nil {
		return "", r.err
	}

	key = r.key(key)

	return redis.String(r.Conn.Do("HGET", key, field))
}

func (r *RedisClient) HLEN(key string) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	key = r.key(key)

	return redis.Int(r.Conn.Do("HLEN", key))
}

func (r *RedisClient) HEXISTS(key, field string) (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	key = r.key(key)

	return redis.Bool(r.Conn.Do("HEXISTS", key, field))
}

func (r *RedisClient) HGETALL(key string) (map[string]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	key = r.key(key)

	return redis.StringMap(r.Conn.Do("HGETALL", key))
}

func (r *RedisClient) INCR(key string) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	key = r.key(key)

	return redis.Int64(r.Conn.Do("INCR", key))
}

func (r *RedisClient) HDEL(key, field string) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	_, err := redis.Int(r.Conn.Do("HDEL", key, field))

	return err
}

func (r *RedisClient) HSCAN(key string, cursor interface{}, optionArgs ...interface{}) (uint64, map[string]string, error) {
	if r.err != nil {
		return 0, nil, r.err
	}

	key = r.key(key)

	args := redis.Args{}.Add(key, cursor).AddFlat(optionArgs)
	result, err := redis.Values(r.Conn.Do("HSCAN", args...))
	if err != nil {
		return 0, nil, err
	}

	newCursor, err := redis.Uint64(result[0], nil)
	if err != nil {
		return 0, nil, err
	}
	data, err := redis.StringMap(result[1], nil)

	return newCursor, data, err
}

func (r *RedisClient) ZADD(key string, score, member interface{}, optionArgs ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	args := redis.Args{}.Add(key).AddFlat(optionArgs)
	_, err := redis.Int(r.Conn.Do("ZADD", args...))
	return err
}

func (r *RedisClient) ZINCRBY(key string, increment, member interface{}) error {
	if r.err != nil {
		return r.err
	}

	key = r.key(key)

	_, err := redis.String(r.Conn.Do("ZINCRBY", key, increment, member))

	return err
}

func (r *RedisClient) Close() {
	if r.Conn != nil {
		r.Conn.Close()
	}
}

func (r *RedisClient) key(key string) string {
	if r.HasPrefix {
		return KeyPrefix + ":" + key

	}

	return key
}