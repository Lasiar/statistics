package db

import (
	"fmt"
	"github.com/kshvakov/clickhouse"
	"github.com/satori/go.uuid"
	"log"
	"statistics/lib"
	"statistics/system"
	"strings"
//	"github.com/go-redis/redis"
)

const dbClickhouseGoodQuery = "INSERT INTO statistics(point_id, played, md5, len) VALUES (?, ?, toFixedString(?, 32),  ?)"

func CheckRedis() bool {
	_, err := lib.RedisStatDB.Ping().Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func SendInfo(ip string, userAgent string, point string) {
	err := lib.RedisIpDB.Set(fmt.Sprint(point, "_ip"), ip, 0).Err()
	if err != nil {
		log.Println(err)
	}
	err = lib.RedisIpDB.Set(fmt.Sprint(point, "_user"), userAgent, 0).Err()
	if err != nil {
		log.Println(err)
	}
}

func SetRedis(statJS lib.StatJS) bool {
	id := uuid.NewV4()
	err := lib.RedisStatDB.Set(fmt.Sprint(id, "_ip:", statJS.Info.Addr, "user_agent:", statJS.Info.Uagent), statJS.Json, 0).Err()
	if err != nil {
		log.Println("Redis set stat:", err)
		return false
	}
	return true
}

func SendToClick(array []lib.ValidJS) error {
	if err := lib.ClickDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			return fmt.Errorf("Error connect tp GoodClick: ", err)
		}
	}
	var (
		tx, _ = lib.ClickDB.Begin()
	)
	stmt, err := tx.Prepare(dbClickhouseGoodQuery)
	if err != nil {
		log.Println(err)
	}
	for _, query := range array {
		if _, err := stmt.Exec(query.Point,
			query.Datetime,
			query.Md5,
			query.Len); err != nil {
			log.Println(err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetStatFromRedis(toParse chan []lib.StatJS) {
	var statArray []lib.StatJS
	var stat lib.StatJS

	KeyDB, err := lib.RedisStatDB.Keys("*ip:*").Result()
	if err != nil {
		log.Println(err)
	}
	if len(KeyDB) == 0 {
		return
	}
	valArr, err := lib.RedisStatDB.MGet(KeyDB...).Result()
	if err != nil {
		log.Println(err)
		return
	}
	for i, val := range valArr {
		d := strings.Index(KeyDB[i], "ip:")
		u := strings.Index(KeyDB[i], "user_agent")
		stat.Info.Addr = KeyDB[i][d+3 : u]
		stat.Info.Uagent = KeyDB[i][u+11:]
		stat.Json, err = system.CheckString(val)
		if err != nil {
			log.Println(err)
			return
		}
		statArray = append(statArray, stat)
	}
	lib.RedisStatDB.Del(KeyDB...).Err()
	if err != nil {
		log.Println(err)
	}
	toParse <- statArray
}
