package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Liste des serveurs NTP à vérifier
var ntpServers = []string{
	"time.google.com",
	"time.cloudflare.com",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
	"fake.ntp.org",
}

var unreachableServers = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ntp_unreachable_servers_total",
		Help: "Total number of times NTP servers that couldn't be reached",
	},
	[]string{"server"},
)

func init() {
	prometheus.MustRegister(unreachableServers)
}

func checkNTPServers(sugar *zap.SugaredLogger) {
	var wg sync.WaitGroup

	for _, server := range ntpServers {
		wg.Add(1)

		go func(server string) {
			defer wg.Done()

			ntpTime, err := ntp.Time(server)
			if err != nil {
				sugar.Infof("Je ne peux pas rejoindre: %s", server)
				unreachableServers.WithLabelValues(server).Inc() // Increment the counter with server label
				return
			}
			localTime := time.Now()
			diff := localTime.Sub(ntpTime)
			sugar.Infof("Difference de temps entre la machine local et %s: %v", server, diff)
		}(server)
	}

	wg.Wait()
}

func main() {
	// Declaration du logging
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil) // Expose the metrics on port 8080

	for {
		checkNTPServers(sugar)
		time.Sleep(15 * time.Second)
	}
}
