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
	var err error
	badDBPing := db.CheckBadDB()
	for {
		select {
		case s := <-badJsonChann:
			badJsonArray = append(badJsonArray, s)
		case <-ticker.C:
			switch {
			case len(badJsonArray) == 0:
				continue
			case badDBPing:
				badDBPing, err = db.SendToBadDB(badJsonArray)
				if err != nil {
					log.Println("Send to bad stat: ", err)
					continue
				}
				badJsonArray = nil
			case len(badJsonArray) > 100:
				badJsonArray = nil
			}
		case <-tickerTen.C:
			badDBPing = db.CheckBadDB()
		}
	}
}

func ReceivingStatWorker(ticker *time.Ticker, halfTicker *time.Ticker, stat chan lib.StatJS, forParse chan []lib.StatJS, count chan int) {
	redisPing := db.CheckRedis()
	var err error
	for {
		select {
		case s := <-stat:
			count <- 1
			if redisPing {
				redisPing, err = db.SetRedis(s)
				if err != nil {
					log.Println("redis stat set: ", err)
				}
			} else {
				sArr := []lib.StatJS{s}
				forParse <- sArr
			}
		case <-halfTicker.C:
			if redisPing {
				err = db.GetStatFromRedis(forParse)
				if err != nil {
					log.Println("redis get stat:", err)
				}
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
			switch {
			case len(infoPointArray) == 0:
				continue
			case redisIpPing:
				redisIpPing, err = db.SendInfo(infoPointArray)
				if err != nil {
					log.Println("Error send ip add: ", err)
					continue
				}
				infoPointArray = nil
			case len(infoPointArray) > 100:
				infoPointArray = nil
			}
		case <-tenTicker.C:
			redisIpPing = db.CheckIpRedis()
		}
	}
}

func ParserWorker(ticker *time.Ticker, statFromRedis chan []lib.StatJS, sendInfoPoint chan lib.InfoPoint, sendBadDB chan lib.BadJS, validJS chan []lib.ValidJS) {
	var arrayValidJS []lib.ValidJS
	for {
		select {
		case <-ticker.C:
			if len(arrayValidJS) == 0 {
				continue
			}
			validJS <- arrayValidJS
			arrayValidJS = nil
		default:
			select {
			case s := <-statFromRedis:
				arrayValidJS = append(arrayValidJS, parser.Parse(s, sendInfoPoint, sendBadDB)...)
			case <-ticker.C:

				if len(arrayValidJS) != 0 {
					validJS <- arrayValidJS
					arrayValidJS = nil
				}
			}
		}
	}
}

func SendClick(ticker *time.Ticker, validJS chan []lib.ValidJS) {
	var arrayValidJS []lib.ValidJS
	clickPing := db.CheckClick()

	for {
		select {
		case <-ticker.C:
			if len(arrayValidJS) == 0 {
				continue
			}

			if clickPing {
				err := db.SendToClick(arrayValidJS)
				if err != nil {
					log.Println("Send Clickhouse: ", err)
					continue
				}
			}
			arrayValidJS = nil
			if len(arrayValidJS) > 950 {
				arrayValidJS = nil
			}
			continue
		default:
			select {
			case <-ticker.C:

				if len(arrayValidJS) == 0 {
					continue
				}
				if clickPing {
					err := db.SendToClick(arrayValidJS)
					if err != nil {
						log.Println("Send Clickhouse: ", err)
						continue
					}
				}
				arrayValidJS = nil
				if len(arrayValidJS) > 950 {
					arrayValidJS = nil
				}
			case a := <-validJS:
				arrayValidJS = append(arrayValidJS, a...)
			}
		}
	}
}
