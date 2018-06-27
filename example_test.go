package sadd

import (
	"fmt"
	"log"
)

func ExampleParseQuery() {
	query := ":6379,:3000-:3003,localhost:3000-:3003,192.168.1.126:3000-:3003,192.168.1.254:3000-192.168.2.1:3001"
	addresses, err := ParseQuery(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addresses {
		fmt.Println(address)
	}
	// outputs:
	//  :6379
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
	//  192.168.1.254:3000
	//  192.168.1.254:3001
	//  192.168.1.255:3000
	//  192.168.1.255:3001
	//  192.168.2.0:3000
	//  192.168.2.0:3001
	//  192.168.2.1:3000
	//  192.168.2.1:3001
}
