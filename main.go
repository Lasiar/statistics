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
	everTenSecond1 := time.NewTicker(10 * time.Second)
	everTenSecond2 := time.NewTicker(10 * time.Second)
	everSecond := time.NewTicker(1 * time.Second)
	everSecond1 := time.NewTicker(1 * time.Second)
	everHalfSecond := time.NewTicker(1 * time.Second)
	everHalfSecond2 := time.NewTicker(1 * time.Second)
	everSecondForClick := time.NewTicker(1 * time.Second)

	stat := make(chan lib.StatJS)
	statFromRedis := make(chan []lib.StatJS)
	sendInfoPoint := make(chan lib.InfoPoint)
	sendBadDB := make(chan lib.BadJS)
	validJS := make(chan []lib.ValidJS)

	go core.SendClick(everSecondForClick, validJS)
	go core.SendRedisIp(everHalfSecond, everTenSecond, sendInfoPoint)
	go core.ReceivingStatWorker(everTenSecond1, everHalfSecond2, stat, statFromRedis)
	go core.ParserWorker(everSecond, statFromRedis, sendInfoPoint, sendBadDB, validJS)
	go core.SendBadJson(everSecond1, everTenSecond2, sendBadDB)


	HandleWeb := web.Web(stat)

	http.HandleFunc("/gateway/statistics/create", HandleWeb)
	http.ListenAndServe(lib.Config.Port, nil)
	fmt.Println("Listen: ", lib.Config.Port)
}
