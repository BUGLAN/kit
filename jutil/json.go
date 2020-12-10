package jutil

import (
	"encoding/json"
	"log"
	"os"
    "fmt"

	"github.com/k0kubun/pp"
)

func OutJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		log.Printf("OutJSON MarshalIndent fail, err is %v \n", err)
		return
	}
	_, _ = os.Stdout.Write(b)
    fmt.Println()
}

func OutStruct(data interface{}) {
	_, err := pp.Print(data)
	if err != nil {
		log.Println(err)
	}
}
