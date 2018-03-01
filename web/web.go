package web

import (
	"fmt"
	"net"
	"net/http"
	"statistics/lib"
	"strings"
)

func Web(stat chan lib.StatJS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statJS lib.StatJS
		fmt.Fprint(w, `{"success":true}`)
		statJS.Json = r.PostFormValue("data")
		statJS.Info.Addr = getRealAddr(r)
		statJS.Info.Uagent = r.UserAgent()
		stat <- statJS
	}
}

func Count(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, lib.Count)
}

func getRealAddr(r *http.Request) string {
	remoteIP := ""
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}
