package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/bigwhite/tcp-server-demo2/frame"
	"github.com/bigwhite/tcp-server-demo2/packet"
	"github.com/lucasepe/codename"
)

func startNewConn() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	frameCodec := frame.NewMyFrameCodec()
	var counter int

	go func() {
		for {
			// handle ack
			// read from the connection
			ackFramePayLoad, err := frameCodec.Decode(conn)
			if err != nil {
				panic(err)
			}

			p, err := packet.Decode(ackFramePayLoad)

			_, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("not submitack")
			}
			//fmt.Printf("the result of submit ack[%s] is %d\n", submitAck.ID, submitAck.Result)
		}
	}()

	for {
		// send submit
		counter++
		id := fmt.Sprintf("%08d", counter) // 8 byte string
		payload := codename.Generate(rng, 4)
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}
		//fmt.Printf("send submit id = %s, payload=%s\n", s.ID, s.Payload)

		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}

		err = frameCodec.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 100)

	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			startNewConn()
			wg.Done()
		}()
	}
	wg.Wait()
}
