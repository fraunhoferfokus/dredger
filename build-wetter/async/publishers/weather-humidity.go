package publishers

import (
	"asyncservice/entities"
	"encoding/json"
	"log" //noch an core logger anpassen
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

func publishHumidityReading() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()


	//err = nc.Publish("weather-humidity", data)

    //Add your business logic here
}