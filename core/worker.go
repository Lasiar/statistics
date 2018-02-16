package core

import (
	"log"
	"statistics/db"
	"statistics/lib"
	"statistics/parser"
	"time"
)

func SendBadJson(ticker *time.Ticker, tickerTen *time.Ticker, badJsonChann chan lib.BadJS) {
	var badJsonArray []lib.BadJS
	badDBPing := db.CheckBadDB()
	for {
		select {
		case s := <-badJsonChann:
			badJsonArray = append(badJsonArray, s)
		case <-ticker.C:
			if len(badJsonArray) == 0 {
				continue
			}
			if badDBPing {
				badDBPing = db.SendToBadDB(badJsonArray)
				badJsonArray = nil
			}
		case <-tickerTen.C:
			badDBPing = db.CheckBadDB()
		}
	}
}

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

func SendRedisIp(ticker *time.Ticker, tenTicker *time.Ticker, infoPoint chan lib.InfoPoint) {
	var infoPointArray []lib.InfoPoint
	var err error
	redisIpPing := db.CheckIpRedis()
	for {
		select {
		case s := <-infoPoint:
			infoPointArray = append(infoPointArray, s)
		case <-ticker.C:
			if len(infoPointArray) != 0 && redisIpPing {
				redisIpPing, err = db.SendInfo(infoPointArray)
				if redisIpPing {
					infoPointArray = nil
				}
				if err != nil {
					log.Println("Error send ip add: ", err)
				}
			}
		case <-tenTicker.C:
			redisIpPing = db.CheckIpRedis()
		}
	}
}

func ParserWorker(ticker *time.Ticker, tenTicker *time.Ticker, stat chan lib.StatJS, statFromRedis chan []lib.StatJS, sendInfoPoint chan lib.InfoPoint, sendBadDB chan lib.BadJS) {
	var arrayValidJS []lib.ValidJS
	clickPing := db.CheckClick()
	returnChannel := make(chan []lib.ValidJS)
	for {
		select {
		case s := <-statFromRedis:
			go parser.Parse(s, returnChannel, sendInfoPoint, sendBadDB)
		case s := <-stat:
			go parser.ParserWithoutRedis(s, returnChannel, sendInfoPoint, sendBadDB)
		case r := <-returnChannel:
			arrayValidJS = append(arrayValidJS, r...)
		case <-ticker.C:
			if len(arrayValidJS) != 0 && clickPing {
				if err := db.SendToClick(arrayValidJS); err != nil {
					log.Println("Send Clickhouse", err)
					continue
				}
				arrayValidJS = nil
			}
			if len(arrayValidJS) > 950 {
				arrayValidJS = nil
			}
		case <-tenTicker.C:
			clickPing = db.CheckClick()
		}
	}
}
