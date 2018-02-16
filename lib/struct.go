package lib

type StatJS struct {
	Json string
	Info struct {
		Addr   string
		Uagent string
	}
}

type InfoPoint struct {
	Point  int
	Addr   string
	Uagent string
}

type Configure struct {
	Redis struct {
		Password string `json:"password"`
		Address  string `json:"address"`
	} `json:"redis"`
	RedisIp struct {
		Password string `json:"password"`
		Address  string `json:"address"`
	} `json:"redis_ip"`
	Psql       string `json:"psql"`
	Clickhouse string `json:"clickhouse"`
}

type RawJS struct {
	Point      int             `json:"point"`
	Statistics [][]interface{} `json:"statistics"`
}

type BadJS struct {
	Ip   string
	Json string
}

type ValidJS struct {
	Point    int
	Datetime int64
	Md5      string
	Len      int
}
