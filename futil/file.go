package futil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func WriteFile(filename string, data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("WriteFile json marshal fail, err is ", err)
	}

	err = ioutil.WriteFile(filename, b, 0666)
	if err != nil {
		fmt.Println("WriteFile ioutil writefile fail, err is ", err)
	}
}
