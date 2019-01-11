package utils

import (
	"encoding/json"
	"log"
)

func Dumper(infos ...interface{}) {
	for _, v := range infos {
		j, _ := json.MarshalIndent(v, "", "\t")
		log.Println(string(j))
	}
}
