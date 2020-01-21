package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lew/persona/app/config"
	log "github.com/sirupsen/logrus"
)

// const (
// 	RedisURL            = "redis://127.0.0.1:6379/3"
// 	redisMaxIdle        = 3   //最大空闲连接数
// 	redisIdleTimeoutSec = 240 //最大空闲连接时间
// 	RedisPassword       = ""
// )

var PersonaRD *RDInstal
var games map[float64]*RDInstal = map[float64]*RDInstal{}

type RDInstal struct {
	pool *redis.Pool
}

func InitRedis() {
	PersonaRD = &RDInstal{newRedisPoolByConfig("redis")}
	// 初始化游戏服的各个数据库
	gamesArr := config.Cfg.MustValueArray("redies", "games", ",")
	for _, j := range gamesArr {
		tempPoll := &RDInstal{newRedisPoolByConfig(j)}
		if tempPoll != nil {
			games[config.Cfg.MustFloat64(j, "appId", 200005)] = tempPoll
		} else {
			sectionVal, _ := config.Cfg.GetSection(j)
			log.Fatalf("redis init Fatal cfgName:%s section:%s ", j, sectionVal)
		}

	}
}
func GetRDIns(appId float64) *RDInstal {
	return games[appId]
}
func newRedisPoolByConfig(configName string) *redis.Pool {
	if _, err := config.Cfg.GetValue(configName, "host"); err != nil {
		log.Panicf("newRedisPoolByConfig config:%s have err:%s ", configName, err.Error())
		return nil
	} else {
		return NewRedisPool(
			config.Cfg.MustValue(configName, "host", "redis://127.0.0.1:6379/3"),
			config.Cfg.MustInt(configName, "maxIdle", 3),
			time.Duration(config.Cfg.MustInt(configName, "idleTimeout", 240))*time.Second,
		)
	}

}

// NewRedisPool 返回redis连接池
// func NewRedisPool(redisURL string) *redis.Pool {

// 	return &redis.Pool{
// 		MaxIdle:     3,
// 		IdleTimeout: 240 * time.Second,
// 		Dial: func() (redis.Conn, error) {
// 			c, err := redis.DialURL(redisURL)
// 			if err != nil {
// 				return nil, fmt.Errorf("redis connection error: %s", err)
// 			}
// 			//验证redis密码
// 			if "" != "" {
// 				if _, authErr := c.Do("AUTH", ""); authErr != nil {
// 					return nil, fmt.Errorf("redis auth password error: %s", authErr)
// 				}
// 			}
// 			return c, err
// 		},
// 		TestOnBorrow: func(c redis.Conn, t time.Time) error {
// 			_, err := c.Do("PING")
// 			if err != nil {
// 				return fmt.Errorf("ping redis error: %s", err)
// 			}
// 			return nil
// 		},
// 	}

// }

func NewRedisPool(uri string, maxIdle int, timeout time.Duration) *redis.Pool {
	log.Infof("NewRedisPool uri:%s ", uri)
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			// redis uri方式 redis://:foobared@10.0.50.11:6379/4
			c, err := redis.DialURL(uri)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			return c, err

			// redis 正常方式
			// con, err := redis.Dial("tcp", conf["Host"].(string),
			// 	redis.DialPassword(conf["Password"].(string)),
			// 	redis.DialDatabase(int(conf["Db"].(int64))),
			// 	redis.DialConnectTimeout(timeout*time.Second),
			// 	redis.DialReadTimeout(timeout*time.Second),
			// 	redis.DialWriteTimeout(timeout*time.Second))
			// if err != nil {
			// 	return nil, fmt.Errorf("redis connection error: %s", err)
			// }
			// return con, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}

	// return &redis.Pool{
	// 	MaxIdle:     config.Cfg.MustInt("redise", "maxIdle", 3),
	// 	IdleTimeout: time.Duration(config.Cfg.MustInt("redise", "maxIdle", 240)) * time.Second,
	// 	Dial: func() (redis.Conn, error) {
	// 		c, err := redis.DialURL(redisURL)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("redis connection error: %s", err)
	// 		}
	// 		//验证redis密码
	// 		if config.Cfg.MustValue("redise", "password", "") != "" {
	// 			if _, authErr := c.Do("AUTH", config.Cfg.MustValue("redise", "password", "")); authErr != nil {
	// 				return nil, fmt.Errorf("redis auth password error: %s", authErr)
	// 			}
	// 		}
	// 		return c, err
	// 	},
	// 	TestOnBorrow: func(c redis.Conn, t time.Time) error {
	// 		_, err := c.Do("PING")
	// 		if err != nil {
	// 			return fmt.Errorf("ping redis error: %s", err)
	// 		}
	// 		return nil
	// 	},
	// }

}

func (rdInstal *RDInstal) Set(k, v string) {
	// c := NewRedisPool(RedisURL).Get()
	c := rdInstal.pool.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		log.Warning(err.Error())
	}
	fmt.Println("set k v", k, v)
}

func (rdInstal *RDInstal) GetStringValue(k string) string {
	// c := NewRedisPool(RedisURL).Get()
	c := rdInstal.pool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		log.Warning(err.Error())
		return ""
	}
	return username
}

func (rdInstal *RDInstal) SetKeyExpire(k string, ex int) {
	// c := NewRedisPool(RedisURL).Get()
	c := rdInstal.pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		log.Warning(err.Error())
	}
}

func (rdInstal *RDInstal) CheckKey(k string) bool {
	c := rdInstal.pool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		log.Warning(err.Error())
		return false
	} else {
		return exist
	}
}

func (rdInstal *RDInstal) DelKey(k string) error {
	c := rdInstal.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		log.Warning(err.Error())
		return err
	}
	return nil
}

func (rdInstal *RDInstal) SetJson(k string, data interface{}) error {
	c := rdInstal.pool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func (rdInstal *RDInstal) getJsonByte(k string) ([]byte, error) {
	c := rdInstal.pool.Get()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}
	return jsonGet, nil
}

func (rdInstal *RDInstal) GetHGET(hash string, k string) ([]byte, error) {
	c := rdInstal.pool.Get()
	jsonGet, err := redis.Bytes(c.Do("HGET", hash, k))
	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}
	return jsonGet, nil
}
func (rdInstal *RDInstal) SetHSET(hash string, k string, data interface{}) error {
	c := rdInstal.pool.Get()
	defer c.Close()
	value, err := json.Marshal(data)
	if err != nil {
		log.Warning(err.Error())
		return err
	}
	_, errDo := c.Do("HSET", hash, k, value)
	if errDo != nil {
		log.Warning(errDo.Error())
		return errDo
	}
	return nil
}

// func (rdInstal *RDInstal) HGETALL(hash string) (map[string]string, error) {
// 	c := rdInstal.pool.Get()
// 	jsonGet, err := redis.StringMap(c.Do("HGETALL", hash))
// 	if err != nil {
// 		log.Warning(err.Error())
// 		return nil, err
// 	}
// 	return jsonGet, nil
// }

func (rdInstal *RDInstal) HGETALL(hash string) (map[string]string, error) {
	c := rdInstal.pool.Get()
	jsonGet, err := redis.StringMap(c.Do("HGETALL", hash))
	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}
	return jsonGet, nil
}

// 设置过期时间 key:键名 expire:过期时间(秒) replace:是否替换
func (rdInstal *RDInstal) SetExpire(key string, expire int64, replace bool) (err error) {
	c := rdInstal.pool.Get()
	defer c.Close()
	// 判断过期时间
	if replace {
		_, err = c.Do("EXPIRE", key, expire)
	} else {
		ttlV, _ := c.Do("TTL", key)
		if ttlV.(int64) == -1 {
			_, err = c.Do("EXPIRE", key, expire)
		}
	}
	return
}
