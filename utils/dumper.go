package utils

import (
	"encoding/json"
	"log"
)

var DumperFlag = true

func Dumper(infos ...interface{}) {
	if DumperFlag {
		j, _ := json.MarshalIndent(infos, "", "\t")
		log.Println(string(j))
	}
}
