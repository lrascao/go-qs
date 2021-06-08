package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/lrascao/go-qs/internal/config"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	consumers cmap.ConcurrentMap
)

func Producer() {

	//
	for {
		conn, err := net.Dial("tcp", config.Cfg.TopEndpoint)
		if err != nil {
			log.Fatal("failed to connect to top endpoint")
		}
		defer conn.Close()
		fmt.Printf("connected to %s\n", config.Cfg.TopEndpoint)

		data := make([]byte, 256)
		for {
			_, err := conn.Read(data)
			if err != nil {
				// reconnect
				break
			}

			// forward data to bottom layer consumers
			for consumer := range consumers.IterBuffered() {
				var c chan []byte
				c = consumer.Val.(chan []byte)
				c <- data
			}
		}
	}
}

func handleConsumer(c chan []byte, conn net.Conn) {
	// block on channel receive
	var data []byte

	select {
	case c <- data:
		// this write might take a long time
		// easy way out is just spawning a goroutine to deal with the delay,
		// the proper way would be having some sort of message queue (similar to Erlang's)
		// that we could manipulate and take informed decisions on
		go func() {
			conn.Write(data)
		}()
	}
}

func main() {
	consumers = cmap.New()

	// kick off a goroutine that will connect to the top layer,
	// get data from it and forward it to all consumers that
	// accepted connections from the bottom layer
	go Producer()

	listener, err := net.Listen("tcp", ":40000")
	if err != nil {
		log.Fatal("unable to listen")
	}

	n_consumers := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept")
		}

		// register the consumer channel
		c := make(chan []byte)
		consumers.Set(strconv.Itoa(n_consumers), c)

		go handleConsumer(c, conn)
		n_consumers++
	}
}
