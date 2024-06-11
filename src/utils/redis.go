package utils

import (
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/jianshao/poker_counter/src/config"
)

var (
	gRedisConn redis.Conn = nil
)

func GetRedisConn() redis.Conn {
	if gRedisConn == nil {
		conn, err := redis.Dial("tcp", config.REDIS_ADDR)
		if err == nil {
			gRedisConn = conn
		}
	}
	return gRedisConn
}

func closeRedis() {
	if gRedisConn != nil {
		gRedisConn.Close()
		gRedisConn = nil
	}
}

func GetString(key string) (string, error) {
	conn := GetRedisConn()
	return redis.String(conn.Do("GET", key))
}

func SetString(key, value string, timeout int) error {
	conn := GetRedisConn()
	err := errors.New("")
	if timeout == 0 {
		_, err = conn.Do("SET", key, value)
	} else {
		_, err = conn.Do("SET", key, value, "EX", timeout)
	}

	return err
}

func GetInt(key string) (int, error) {
	conn := GetRedisConn()
	return redis.Int(conn.Do("GET", key))
}

func Inc(key string) (int, error) {
	conn := GetRedisConn()
	return redis.Int(conn.Do("INCR", key))
}

func SetInt(key string, value int) error {
	conn := GetRedisConn()
	_, err := conn.Do("SET", key, value)
	return err
}
