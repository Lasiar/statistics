package main

import (
	"fmt"
	"net/http"
	"statistics/configure"
	"statistics/core"
	"statistics/db"
	"statistics/lib"
	"statistics/system"
	"statistics/web"
	"time"
)

func init() {
	configure.Config()
	db.NewRedisIp()
	db.NewClick()
	db.NewPostSql()
	system.Exit()
	system.GenUUID()
	db.NewRedisStat()

}

func main() {
	everTenSecond := time.NewTicker(10 * time.Second)
	everSecond := time.NewTicker(1 * time.Second)
	everHalfSecond := time.NewTicker(500 * time.Millisecond)

	stat := make(chan lib.StatJS)
	statFromRedis := make(chan []lib.StatJS)
	sendInfoPoint := make(chan lib.InfoPoint)
	sendBadDB := make(chan lib.BadJS)

	go core.SendRedisIp(everHalfSecond, everTenSecond, sendInfoPoint)
	go core.ReceivingStatWorker(everTenSecond, everHalfSecond, stat, statFromRedis)
	go core.ParserWorker(everSecond, everTenSecond, statFromRedis, sendInfoPoint, sendBadDB)
	go core.SendBadJson(everSecond, everTenSecond, sendBadDB)

	HandleWeb := web.Web(stat)

	http.HandleFunc("/gateway/statistics/create", HandleWeb)
	http.ListenAndServe(lib.Config.Port, nil)
	fmt.Println("Listen: ", lib.Config.Port)
}
