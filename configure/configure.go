package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"statistics/lib"
)

func Config() {
	file, err := ioutil.ReadFile("config/config")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &lib.Config)
	if err != nil {
		fmt.Println("Unmarshal config", err)
	}
}
