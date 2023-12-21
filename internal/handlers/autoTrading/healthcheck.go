package autoTrading

import (
	"log"
	"time"
)

func HealthCheck() {
	healthCheckTicker := time.NewTicker(time.Second)

	for {
		select {
		case <-healthCheckTicker.C:
			log.Print("Auto Trading going well")
			healthCheckTicker = time.NewTicker(time.Hour)
		}
	}

}
