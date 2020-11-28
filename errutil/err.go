package errutil

import "fmt"

func CheckErrWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("check err fail, err is", err)
	}
}
