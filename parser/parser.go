package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"statistics/lib"
	"statistics/system"
	"strconv"
	"strings"
)

func Parse(statArray []lib.StatJS, statChannel chan []lib.ValidJS, sendInfoPoint chan lib.InfoPoint, sendBadDB chan lib.BadJS) {
	for _, stat := range statArray {
		js, err := unmarshalJS(stat.Json)
		if err != nil {
			log.Println("Error: ", " ip addr: ", stat.Info.Addr, " user info: ", stat.Info.Uagent, " json: ", stat.Json, err)
			sendBadDB <- system.MakeBadJS(stat)
			return
		}

		sendInfoPoint <- system.MakeInfoPoint(js, stat)

		err = validInterfaceJS(js.Statistics)
		if err != nil {
			log.Println("Error: ", " ip addr: ", stat.Info.Addr, " user info: ", stat.Info.Uagent, " json: ", stat.Json, err)
			sendBadDB <- system.MakeBadJS(stat)
			return
		}
		readyJs, err := changeType(js)
		if err != nil {
			log.Println("Error: ", " ip addr: ", stat.Info.Addr, " user info: ", stat.Info.Uagent, " json: ", stat.Json, err)
			sendBadDB <- system.MakeBadJS(stat)
			return
		}
		statChannel <- readyJs
	}
}

func unmarshalJS(js string) (lib.RawJS, error) {
	var rawJson lib.RawJS
	jsBool := []byte(js)
	err := json.Unmarshal([]byte(jsBool), &rawJson)
	if err != nil {
		return rawJson, fmt.Errorf("Error unmarshal: %v", err)
	}
	return rawJson, nil
}

func validInterfaceJS(JSArray [][]interface{}) error {
	if len(JSArray) == 0 {
		return fmt.Errorf("%v", "Array in json empty")
	}
	for _, firstArray := range JSArray {
		for i, secondArray := range firstArray {
			switch t := secondArray.(type) {
			case float64:
				if i == 0 || i == 1 {
					return fmt.Errorf("Error: invalid json: want float 64, have %v", t)
				}
			case string:
				if i == 2 {
					return fmt.Errorf("Error: invalid json: want string, have %v", t)
				} else {
					if i == 1 {
						if strings.Contains(t, " ") {
							return fmt.Errorf("Error: invalid json: md5 has space")
						}
						if l := len(t); l != 32 {
							return fmt.Errorf("Error: invalid json: want md5 lenght 32, have %v", t)
						}
					} else {
						if strings.Contains(t, " ") {
							return fmt.Errorf("Error: invalid json: datetime has space")
						}
					}
				}
			default:
				return fmt.Errorf("Error: unknow type %v", t)
			}
		}
	}
	return nil
}

func changeType(rawJson lib.RawJS) ([]lib.ValidJS, error) {
	var err error
	LenQuery := len(rawJson.Statistics)
	query := make([]lib.ValidJS, LenQuery, LenQuery)
	for p, first := range rawJson.Statistics {
		query[p].Point = rawJson.Point
		for i, second := range first {
			switch i {
			case 0:
				query[p].Datetime, err = strconv.ParseInt(second.(string), 10, 64)
				if err != nil {
					return query, err
				}
			case 1:
				query[p].Md5 = second.(string)
			case 2:
				query[p].Len = int(second.(float64))
			}
		}
	}
	return query, nil
}
