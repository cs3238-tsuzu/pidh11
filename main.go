package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/d2r2/go-dht"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	temperatureGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tsuzu_room_temperature",
		Help: "The current temperature in Tsuzu's room",
	})
	humidityGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tsuzu_room_humidity",
		Help: "The current humidity in Tsuzu's room",
	})
)

func main() {

	go func() {
		for {
			temperature, humidity, _, err :=
				dht.ReadDHTxxWithRetry(dht.DHT11, 14, false, 20)
			if err != nil {
				log.Fatal(err)
			}

			if err != nil {
				log.Fatal(err)
			}

			temperatureGauge.Set(float64(temperature))
			humidityGauge.Set(float64(humidity))

			fmt.Printf("temperature: %f, humidity: %f\n", temperature, humidity)

			time.Sleep(3 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	server := http.Server{
		Addr: ":2112",
	}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	server.ListenAndServe()
}
