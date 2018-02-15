package core

import (
	"fmt"
	"log"
	"statistics/db"
	"statistics/lib"
	"statistics/parser"
	"time"
)




func ReceivingStatWorker(ticker *time.Ticker, halfTicker *time.Ticker, stat chan lib.StatJS, sendParse chan lib.StatJS, forParse chan []lib.StatJS) {
	redisPing := db.CheckRedis()
	for {
		select {
		case s := <-stat:
			if redisPing {
				redisPing = db.SetRedis(s)
			} else {
				sendParse <- s
			}
		case <-halfTicker.C:
			if redisPing {
				db.GetStatFromRedis(forParse)
			}
		case <-ticker.C:
			redisPing = db.CheckRedis()

		}
	}
}

func ParserWorker(ticker *time.Ticker, stat chan lib.StatJS, statFromRedis chan []lib.StatJS) {
	var arrayValidJS []lib.ValidJS
	returnChannel := make(chan []lib.ValidJS)
	for {
		select {
		case s := <-statFromRedis:
			go parser.Parse(s, returnChannel)
		case s := <-stat:
			go parser.ParserWithoutRedis(s, returnChannel)
		case r := <-returnChannel:
			arrayValidJS = append(arrayValidJS, r...)
			fmt.Println(len(arrayValidJS))
		case <-ticker.C:
			if len(arrayValidJS) != 0 {
				fmt.Println(arrayValidJS)
				if err := db.SendToClick(arrayValidJS); err != nil {
					log.Println("Send Clickhouse", err)
					continue
				}
				arrayValidJS = nil
			}
		}
	}
}
