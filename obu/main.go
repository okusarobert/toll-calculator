package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/okusarobert/toll-calculator/types"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second

func sendOBUData(conn *websocket.Conn, data types.OBUData) error {
	return conn.WriteJSON(data)
}

func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLocation()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			if err := sendOBUData(conn, data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}

func init() {

	rand.Seed(time.Now().UnixNano())
}
