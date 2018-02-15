package main

import (
	"net/http"
	"statistics/configure"
	"statistics/core"
	"statistics/db"
	"statistics/lib"
	"statistics/web"
	"time"
)

func init() {
	configure.Config()
	db.NewRedisStat()
	db.NewRedisIp()
	db.NewClick()
	db.NewPostSql()
}

func main() {
	everTenSecond := time.NewTicker(10 * time.Second)
	everSecond := time.NewTicker(1 * time.Second)
	everHalfSecond := time.NewTicker(500 * time.Millisecond)

	stat := make(chan lib.StatJS)
	sendInParse := make(chan lib.StatJS)
	statFromRedis := make(chan []lib.StatJS)

	go core.ReceivingStatWorker(everTenSecond, everHalfSecond, stat, sendInParse, statFromRedis)
	go core.ParserWorker(everSecond, sendInParse, statFromRedis)

	HandleWeb := web.Web(stat)

	http.HandleFunc("/gateway/statistics/create", HandleWeb)
	http.ListenAndServe(":8080", nil)
}
