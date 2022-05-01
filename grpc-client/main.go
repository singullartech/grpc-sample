package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-stream/pkg/grpc"
	"log"
	"sync"
	"time"

	gogrpc "google.golang.org/grpc"
)

var (
	concurrency int
	servers     int
	messages    int
)

func init() {
	flag.IntVar(&concurrency, "concurrency", 1, "number of client routines")
	flag.IntVar(&servers, "servers", 1, "number of grpc servers to use")
	flag.IntVar(&messages, "messages", 1, "number of message to be send")
	flag.Parse()
}

func main() {
	log.Println("-----------------------------------------")
	log.Println("Starting grpc client")
	var conns []*gogrpc.ClientConn
	var wg sync.WaitGroup

	wg.Add(concurrency)

	log.Println("----------------Routines-----------------")
	start := time.Now()
	for i, su := 1, 1; i <= concurrency; i, su = i+1, su+1 {
		var conn *gogrpc.ClientConn

		if su > servers {
			su = 1
		}

		conn, err := openClientConnection(i, su)
		if err != nil {
			log.Fatalln("error openning grpc client connection", err)
		}

		conns = append(conns, conn)
	}

	log.Println("----------------Responses-----------------")
	for index, conn := range conns {
		go send(conn, index+1, &wg, int(messages/concurrency))
	}

	wg.Wait()

	for _, conn := range conns {
		defer conn.Close()
	}

	end := time.Now()
	spent := end.Sub(start)

	log.Println("-----------------------------------------")
	log.Println("# Routines", concurrency)
	log.Println("# Servers", servers)
	log.Println("# Messages", messages)
	log.Println("Time spent (seconds): ", spent.Seconds())
	log.Println("-----------------------------------------")
}

func openClientConnection(t int, s int) (*gogrpc.ClientConn, error) {
	log.Println(fmt.Sprintf("routine %03d: grpc-server-%d:50051", t, s))
	return gogrpc.Dial(fmt.Sprintf("grpc-server-%d:50051", s), gogrpc.WithInsecure(), gogrpc.WithBlock())
}

func send(conn *gogrpc.ClientConn, pool int, wg *sync.WaitGroup, messages int) {
	defer wg.Done()

	client := grpc.NewMessageStream(conn, messages)
	if err := client.Send(context.Background(), pool); err != nil {
		log.Println("error sending message: ", err)
	}
}
