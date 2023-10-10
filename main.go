package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

// Liste des serveurs NTP à vérifier
var ntpServers = []string{
	"time.google.com",
	"time.cloudflare.com",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
	"fake.ntp.org",
}

// On prend le temps courant et on ajoute 10 minutes.
const timeDifferenceThreshold = 10 * time.Minute

func main() {
	for _, server := range ntpServers {
		// On recoit le temps de chaque serveurs NTP
		ntpTime, err := ntp.Time(server)
		if err != nil {
			log.Printf("Error fetching time from %s: %v", server, err)
			continue
		}

		// On calcule la difference de temps entre le pod et le serveur NTP
		localTime := time.Now()
		diff := localTime.Sub(ntpTime)

		fmt.Printf("Time difference between local machine and %s: %v\n", server, diff)
	}
}
