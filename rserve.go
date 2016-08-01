package main

import (
	"fmt"
	"io/ioutil"

	"github.com/senseyeio/roger"
)

func main() {
	rClient, err := roger.NewRClient("192.168.64.2", 6311)
	if err != nil {
		fmt.Println("Failed to connect")
		return
	}

	value, err := rClient.Eval("pi")
	if err != nil {
		fmt.Println("Command failed: " + err.Error())
	} else {
		fmt.Println(value) // 3.141592653589793
	}

	content, err := ioutil.ReadFile("fitCalculator.R")
	if err != nil {
		panic(err)
	}

	helloWorld, err := rClient.Eval(string(content) + "calculateFit(\"potato\")")
	if err != nil {
		panic(err)
	}
	fmt.Println(helloWorld) // Hello World

	arrChan := rClient.Evaluate("Sys.sleep(5); c(1,1)")
	arrResponse := <-arrChan
	arr, _ := arrResponse.GetResultObject()
	fmt.Println(arr) // [1, 1]
}
