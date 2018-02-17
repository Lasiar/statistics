package db

import (
	"fmt"
	"github.com/kshvakov/clickhouse"
	"github.com/satori/go.uuid"
	"log"
	"statistics/lib"
	"statistics/system"
	"strings"
)

const (
	dbClickhouseGoodQuery = "INSERT INTO statistics(point_id, played, md5, len) VALUES (?, ?, toFixedString(?, 32),  ?)"
	dbClickhouseBadQuery  = `INSERT INTO badjson(ip, json) VALUES ($1, $2)`
)

func SendInfo(infoPointArr []lib.InfoPoint) (bool, error) {
	for _, infoPoint := range infoPointArr {
		err := lib.RedisIpDB.Set(fmt.Sprint(infoPoint.Point, "_ip"), infoPoint.Addr, 0).Err()
		if err != nil {
			return false, fmt.Errorf("%v %v: ", "Set ip addr", err)
		}
		err = lib.RedisIpDB.Set(fmt.Sprint(infoPoint.Point, "_user"), infoPoint.Uagent, 0).Err()
		if err != nil {
			return false, fmt.Errorf("%v %v: ", "Set uagent", err)
		}
	}
	return true, nil
}

func SetRedis(statJS lib.StatJS) (bool, error) {
	id := uuid.NewV4()
	err := lib.RedisStatDB.Set(fmt.Sprint(id, "_ip:", statJS.Info.Addr, "user_agent:", statJS.Info.Uagent), statJS.Json, 0).Err()
	if err != nil {
		return false, fmt.Errorf("%v %v: ", "Set stat", err)
	}
	return true, nil
}

func SendToBadDB(badJsons []lib.BadJS) bool {
	if err := lib.PsqlDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return false
		} else {
			return false
		}
	}
	var (
		tx, _ = lib.PsqlDB.Begin()
	)
	stmt, err := tx.Prepare(dbClickhouseBadQuery)
	if err != nil {
		log.Println(err)
	}
	for _, query := range badJsons {
		if _, err := stmt.Exec(query.Ip,
			query.Json); err != nil {
			log.Println(err)
			return false
		}
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
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

func GetStatFromRedis(toParse chan []lib.StatJS) error {
	var statArray []lib.StatJS
	var stat lib.StatJS

	KeyDB, err := lib.RedisStatDB.Keys("*ip:*").Result()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if len(KeyDB) == 0 {
		return nil
	}
	valArr, err := lib.RedisStatDB.MGet(KeyDB...).Result()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	for i, val := range valArr {
		d := strings.Index(KeyDB[i], "ip:")
		u := strings.Index(KeyDB[i], "user_agent")
		stat.Info.Addr = KeyDB[i][d+3 : u]
		stat.Info.Uagent = KeyDB[i][u+11:]
		stat.Json, err = system.CheckString(val)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		statArray = append(statArray, stat)
	}
	lib.RedisStatDB.Del(KeyDB...).Err()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	toParse <- statArray
	return nil
}

/*
func SendBadStatistic(js []lib.BadJS) {
	for _, jsonRaw := range js {
		jsonStr, _ := json.Marshal(jsonRaw)
		req, err := http.NewRequest("POST", lib., bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
}
*/
