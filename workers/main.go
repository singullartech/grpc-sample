package main

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Body string `json:"body"`
}

const (
	MongoConnection = "mongodb://mongo-server:27017"
	NatsConnection  = "nats://nats-server:4222"
	MessageSubject  = "sample.sub.messages"
	QueueName       = "workers"
)

var (
	nc   *nats.Conn
	mc   *mongo.Client
	coll *mongo.Collection
	i    int = 0
	li   []interface{}
)

func init() {
	var err error

	nc, err = nats.Connect(NatsConnection)
	if err != nil {
		log.Panic("failed to connect to NATS", err)
	}

	mc, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoConnection))
	if err != nil {
		log.Panic(err)
	}

	coll = mc.Database("sample").Collection("messages")
}

func main() {
	defer func() { workerRecover() }()

	go saver()

	if _, err := nc.QueueSubscribe(MessageSubject, QueueName, func(msg *nats.Msg) {
		data := Message{}

		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Panic("failed to unmarshal message: ", err, string(msg.Data))
		}

		li = append(li, data)
		i++
	}); err != nil {
		log.Panic(err)
	}

	runtime.Goexit()
}

func saver() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		if len(li) > 0 {
			i = 0
			la := li
			li = li[:0]

			r, err := coll.InsertMany(context.TODO(), la)
			if err != nil {
				log.Fatalln("failed to insert messages", err)
			}

			log.Println("messages saved: ", len(r.InsertedIDs))
		}
	}
}

func workerRecover() {
	if r := recover(); r != nil {
		if err := mc.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}

		nc.Drain()
		nc.Close()
	}
}
