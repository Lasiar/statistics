package system

import (
	"fmt"
	"statistics/lib"

	"os"
	"os/signal"
	"syscall"
	"bufio"
	"strings"
	"log"
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

func confirm(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}

func Exit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			<-sigs
			if !confirm("\nyou want exit?", 1) {
				fmt.Println("nocd")
				continue
			}
			fmt.Println("finish work")
			lib.RedisIpDB.Close()
			lib.RedisStatDB.Close()
			lib.PsqlDB.Close()
			lib.ClickDB.Close()
			fmt.Println("exit")
			os.Exit(0)
		}
	}()
}