package main

import (
	"client"
	"fmt"
	"log"
	"time"
)

type MetadataPushDelegate struct{}

func (m MetadataPushDelegate) ServiceRegister(input []byte) {
	log.Println(string(input))
}

func (m MetadataPushDelegate) ServiceUnRegister(input []byte) {
	log.Println(string(input))
}

func main() {
	errorChan := make(chan any)
	c, err := client.Setup("aaabbbcccdddeee", "192.168.233.52", 7878, MetadataPushDelegate{})
	if err != nil {
		log.Fatalln(err)
	}
	err = c.RegisterConsoleListener()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for i := 1; i <= 1000; i++ {
			<-time.After(time.Millisecond * 500)
			c.Console(fmt.Sprintf("hahah %d", i))
			log.Printf("times: %d", i)
		}
	}()
	<-errorChan
}
