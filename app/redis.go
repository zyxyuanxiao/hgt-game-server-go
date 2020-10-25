package app

import (
	"encoding/json"
	redigo "github.com/gomodule/redigo/redis"
	"io/ioutil"
	"time"
)
var Redis *RedisConfig

// redis配置
type RedisConfig struct {
	Ip 	       string
	Port       string
	Password   string
	Database   int64
}

func LoadRedis() {
	var err error
	redisJson, err := ioutil.ReadFile( Path + "/config/" + ENV + "/redis.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(redisJson, &Redis)
	if err != nil {
		panic(err)
	}
}

// redis pool
func (redis *RedisConfig) GetRedisByPool() *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   3,//最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", redis.Ip+":"+redis.Port, redigo.DialPassword(redis.Password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}