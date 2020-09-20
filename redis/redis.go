package redis

// If you Confirure redis settings to config.json and call InitializeRedis.
// you are ready to use Redis connections.

import (
	"strconv"
	"time"

	conf "github.com/chixm/servertemplate2/config"
	"github.com/gomodule/redigo/redis"
	logrus "github.com/sirupsen/logrus"
)

var redisConnections map[string]*redis.Pool

var logger *logrus.Entry

// InitializeRedis Make connections to Redis from Config
func InitializeRedis(l *logrus.Entry) {
	logger = l

	config := conf.GetConfig()

	redisConnections = make(map[string]*redis.Pool)

	for _, r := range config.Redis {
		pool := makeConnectionPool(r.Host, r.Port, r.MaxIdle, r.MaxActive)
		redisConnections[r.Id] = pool

		// test connection
		_, err := pool.Get().Do("set", "test", "testvalue")
		if err != nil {
			panic(err)
		}
		logger.Println("Tested Redis Connection of " + r.Host)
		rep, err := redis.String(pool.Get().Do("get", "test"))
		if err != nil {
			panic(err)
		}
		if rep == "testvalue" {
			logger.Println("Tested Redis " + r.Host + " OK.")
		}
	}
}

func makeConnectionPool(host string, port int, maxIdle, maxActive int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host+":"+strconv.Itoa(port)) },
	}
}

// TerminateRedis redis terminater
func TerminateRedis() {
	for _, c := range redisConnections {
		err := c.Close()
		if err != nil {
			logger.Println("Redis Close Error::" + err.Error())
		}
	}
}
