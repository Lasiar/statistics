package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"statistics/lib"
	"os"
	"log"
)

func Config() {
	file, err := ioutil.ReadFile("config/config")
	if err != nil {
		log.Println("Error read config: ", err, "because of this we will create it")
		createConfig()
		return
	}
	err = json.Unmarshal(file, &lib.Config)
	if err != nil {
		fmt.Println("Unmarshal config", err)
	}
}

func createConfig() {
	os.Mkdir("config/", 0777)
	lib.Config.Port = ":8080"
	lib.Config.Psql = "user=psql_dbname=dbname password=qwerty host=127.0.0.1"
	lib.Config.Clickhouse = "tcp://127.0.0.1:9000?database=stat"
	lib.Config.Redis.Password="qwerty"
	lib.Config.Redis.Address="127.0.0.1"
	lib.Config.RedisIp.Password="qwerty"
	lib.Config.RedisIp.Address="127.0.0.1"
	StrBool ,err := json.MarshalIndent(lib.Config, "", "   ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile("config/config", StrBool, 0777)
	if err != nil {
		log.Println(err)
	}
}