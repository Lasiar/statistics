package db

import (
	"fmt"
	"github.com/kshvakov/clickhouse"
	"log"
	"statistics/lib"
)

func CheckRedis() bool {
	_, err := lib.RedisStatDB.Ping().Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckBadDB() bool {
	if err := lib.PsqlDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return false
		} else {
			log.Println("Error send badDB: ", err)
			return false
		}
	}
	return true
}

func CheckClick() bool {
	if err := lib.ClickDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return false
		} else {
			log.Println("Error send click: ", err)
			return false
		}
	}
	return true
}

func CheckIpRedis() bool {
	_, err := lib.RedisIpDB.Ping().Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
