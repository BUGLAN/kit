package util

import (
	"encoding/json"
	"log"
	"os"
)

func OutJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		log.Printf("OutJSON MarshalIndent fail, err is %v \n", err)
		return
	}
	_, _ = os.Stdout.Write(b)
}


