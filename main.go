package main

import (
	"fmt"
	"io/ioutil"
	"pings/p1"
	"strings"
	"time"
)

func main() {
	//	start := time.Now()
	b, err := ioutil.ReadFile("./ip")
	if err != nil {
		fmt.Println("1111111111", err)
	}
	iplist := strings.Split(string(b), "\n")
	fmt.Println("total ip is ", len(iplist))

	for _, j := range iplist {

		go p1.Goping(j)

	}
	//	elapsed := time.Since(start).Seconds()
	//	fmt.Println(elapsed)
	time.Sleep(time.Second * 3)

}
