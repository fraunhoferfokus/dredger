// Edit this file, as it is a specific handler function for your service
package publishers

import (
	"log" //noch an core logger anpassen

	"github.com/nats-io/nats.go"
)

func PublishHumidityReading() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	//err = nc.Publish("weather-humidity", data)

	//Add your business logic here
}
