package server

import (
	"async-service/entities"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)
func subscribeToTemperatureReadings(nc *nats.Conn, ws *websocket.Conn) *nats.Subscription {

	// Subscribe to "weatherTemperatur"
	Sub, err := nc.Subscribe("weatherTemperatur", func(msg *nats.Msg) {
		env := entities.Envelope{
			Type: "weatherTemperatur",
			Data: msg.Data,
		}
		out, err := json.Marshal(env)
		if err != nil {
			log.Println("Error marshaling weatherTemperatur envelope:", err)
			return
		}
		if writeErr := ws.WriteMessage(websocket.TextMessage, out); writeErr != nil {
			log.Println("Error writing humidity to WebSocket:", writeErr)
		}
	})
	if err != nil {
		log.Println("NATS subscribe to weatherTemperatur failed:", err)
		return nil
	}

	return Sub
}

func publishTemperatureReading(nc *nats.Conn, ws *websocket.Conn) *nats.Subscription {

	// Subscribe to "weatherTemperatur"
	Sub, err := nc.Subscribe("weatherTemperatur", func(msg *nats.Msg) {
		env := entities.Envelope{
			Type: "weatherTemperatur",
			Data: msg.Data,
		}
		out, err := json.Marshal(env)
		if err != nil {
			log.Println("Error marshaling weatherTemperatur envelope:", err)
			return
		}
		if writeErr := ws.WriteMessage(websocket.TextMessage, out); writeErr != nil {
			log.Println("Error writing humidity to WebSocket:", writeErr)
		}
	})
	if err != nil {
		log.Println("NATS subscribe to weatherTemperatur failed:", err)
		return nil
	}

	return Sub
}

func subscribeToHumidityReadings(nc *nats.Conn, ws *websocket.Conn) *nats.Subscription {

	// Subscribe to "weatherHumidity"
	Sub, err := nc.Subscribe("weatherHumidity", func(msg *nats.Msg) {
		env := entities.Envelope{
			Type: "weatherHumidity",
			Data: msg.Data,
		}
		out, err := json.Marshal(env)
		if err != nil {
			log.Println("Error marshaling weatherHumidity envelope:", err)
			return
		}
		if writeErr := ws.WriteMessage(websocket.TextMessage, out); writeErr != nil {
			log.Println("Error writing humidity to WebSocket:", writeErr)
		}
	})
	if err != nil {
		log.Println("NATS subscribe to weatherHumidity failed:", err)
		return nil
	}

	return Sub
}

func publishHumidityReading(nc *nats.Conn, ws *websocket.Conn) *nats.Subscription {

	// Subscribe to "weatherHumidity"
	Sub, err := nc.Subscribe("weatherHumidity", func(msg *nats.Msg) {
		env := entities.Envelope{
			Type: "weatherHumidity",
			Data: msg.Data,
		}
		out, err := json.Marshal(env)
		if err != nil {
			log.Println("Error marshaling weatherHumidity envelope:", err)
			return
		}
		if writeErr := ws.WriteMessage(websocket.TextMessage, out); writeErr != nil {
			log.Println("Error writing humidity to WebSocket:", writeErr)
		}
	})
	if err != nil {
		log.Println("NATS subscribe to weatherHumidity failed:", err)
		return nil
	}

	return Sub
}
