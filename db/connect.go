package db

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kshvakov/clickhouse"
	_ "github.com/lib/pq"
	"log"
	"statistics/lib"
)

func NewRedisStat() {
	lib.RedisStatDB = redis.NewClient(&redis.Options{
		Addr:     lib.Config.Redis.Address,
		Password: lib.Config.Redis.Password, // no password set
		DB:       0,                         // use default DB
	})
	_, err := lib.RedisStatDB.Ping().Result()
	if err != nil {
		log.Println(err)
	}
}

func NewRedisIp() {
	lib.RedisIpDB = redis.NewClient(&redis.Options{
		Addr:     lib.Config.RedisIp.Address,
		Password: lib.Config.RedisIp.Password, // no password set
		DB:       0,                           // use default DB
	})
	_, err := lib.RedisStatDB.Ping().Result()
	if err != nil {
		log.Println(err)
	}
}

func NewPostSql() {
	var err error
	lib.PsqlDB, err = sql.Open("postgres", lib.Config.Psql)
	if err != nil {
		log.Fatal(err)
	}
	if err := lib.PsqlDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
}

func NewClick() {
	var err error
	lib.ClickDB, err = sql.Open("clickhouse", lib.Config.Clickhouse)
	if err != nil {
		log.Fatal(err)
	}
	if err := lib.ClickDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
}
