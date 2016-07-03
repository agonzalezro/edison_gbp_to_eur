package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

const (
	exchangeRateServiceURL = "http://api.fixer.io/latest?symbols=EUR,GBP"
	maxBackoff             = 60 // seconds
)

func getRate() float32 {
	var (
		resp *http.Response
		err  error
	)

	// Yep, you are right, I probably deserve what you are thinking for this backoff :)
	if resp, err = http.Get(exchangeRateServiceURL); err != nil {
		backoff := 1
		ticker := time.NewTicker(time.Duration(backoff) * time.Second)

	LOOP:
		for {
			select {
			case <-ticker.C:
				if resp, err = http.Get(exchangeRateServiceURL); err == nil {
					break LOOP
				}
				backoff = backoff * 2
				ticker = time.NewTicker(time.Duration(backoff) * time.Second)
			case <-time.After(maxBackoff * time.Second):
				// Nothing more that we can do, return an "error"
				return 0.0
			}
		}
	}
	defer resp.Body.Close()

	type Payload struct {
		Rates map[string]float32
	}
	var payload Payload
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return 0.0
	}
	return 1 / payload.Rates["GBP"]
}

func main() {
	gbot := gobot.NewGobot()

	board := edison.NewEdisonAdaptor("edison")
	screen := i2c.NewGroveLcdDriver(board, "screen")

	work := func() {
		rate := getRate()

		screen.Write(fmt.Sprintf("GBP to EUR\n%.3f", rate))

		switch {
		case rate < 1.20:
			screen.SetRGB(255, 0, 0)
		case rate < 1.30:
			screen.SetRGB(0, 0, 255)
		default:
			screen.SetRGB(0, 255, 0)
		}

		<-time.After(3600 * time.Second)
	}

	robot := gobot.NewRobot("edison GBP to EUR",
		[]gobot.Connection{board},
		[]gobot.Device{screen},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
