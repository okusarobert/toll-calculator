package main

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/okusarobert/toll-calculator/types"
	"golang.org/x/exp/rand"
)

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch    chan types.OBUData
	conn     *websocket.Conn
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgch:    make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU client connected")
	for {
		var data types.OBUData
		reqID := rand.Intn(math.MaxInt)
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		data.RequestID = string(reqID)
		dr.producer.ProduceData(data)
	}
}
