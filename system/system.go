package system

import (
	"fmt"
	"statistics/lib"
)

func CheckString(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	default:
		return "", fmt.Errorf("some errors", v)
	}
}

func MakeInfoPoint(js lib.RawJS, statJS lib.StatJS) lib.InfoPoint {
	var inf lib.InfoPoint
	inf.Point = js.Point
	inf.Addr = statJS.Info.Addr
	inf.Uagent = statJS.Info.Uagent
	return inf
}

func MakeBadJS(stat lib.StatJS) lib.BadJS {
	var bad lib.BadJS
	bad.Json = stat.Json
	bad.Ip = stat.Info.Addr
	return bad
}