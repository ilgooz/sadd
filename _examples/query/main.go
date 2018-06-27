package main

import (
	"fmt"
	"log"

	"github.com/ilgooz/sadd"
)

func main() {
	query := ":6379,:3000-:3003,localhost:3000-:3003,192.168.1.126:3000-:3003,192.168.1.254:3000-192.168.2.1:3001"
	addresses, err := sadd.ParseQuery(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addresses {
		fmt.Println(address)
	}
	// outputs:
	// 	:6379
	//  :3000
	//  :3001
	//  :3002
	//  :3003
	//  localhost:3000
	//  localhost:3001
	//  localhost:3002
	//  localhost:3003
	//  192.168.1.126:3000
	//  192.168.1.126:3001
	//  192.168.1.126:3002
	//  192.168.1.126:3003
}
