package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type Move struct {
	date    time.Time
	move    float64
	percent float64
}

type opts struct {
	Apca_api_key_id     string
	Apca_api_secret_key string
	Apca_base_url       string
}

var (
	symbol = flag.String("s", "", "Symbol to look up")
	years  = flag.Int("y", 1, "Years to go back")
	Opts   opts
)

// computeADR takes a slice of any number of bars and returns the average daily range.
func computeADR(bars []marketdata.Bar) float64 {

	var sum float64 = 0
	for _, bar := range bars {
		sum += bar.High - bar.Low
	}
	return sum / float64(len(bars))
}

// getBars requests from the Alpaca API the specified number of years worth of daily bars
// for the specified stock symbol.
func getBars(symbol string, years int) []marketdata.Bar {
	fmt.Printf("Getting %d years of ranges for %v.\n", years, symbol)
	dataClient := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    Opts.Apca_api_key_id,
		ApiSecret: Opts.Apca_api_secret_key,
	})

	now := time.Now()
	// Subtracting 1 day from Now because the free alpaca API doesn't let me pull the last 15 minutes.
	yesterday := now.AddDate(0, 0, -1)
	yearAgo := now.AddDate(-years, 0, 0)
	bars, err := dataClient.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame:  marketdata.NewTimeFrame(1, marketdata.Day),
		Start:      yearAgo,
		End:        yesterday,
		Adjustment: marketdata.Split,
	})
	if err != nil {
		log.Fatalf("Failed to fetch bars for %v %v", symbol, err)
	}
	return bars
}

// getRanges computes the dollar and percentage moves for each bar.
// It then returns an array of Move structs in the order the bars were received.
func getRanges(bars []marketdata.Bar) []Move {
	ranges := make([]Move, 0, len(bars))
	for _, bar := range bars {
		ranges = append(ranges, Move{
			bar.Timestamp,
			bar.Close - bar.Open,
			(bar.Close - bar.Open) / bar.Close * 100,
		})
	}
	return ranges
}

func main() {
	data, _ := ioutil.ReadFile(".env.json")
	err := json.Unmarshal(data, &Opts)
	if err != nil {
		log.Fatal("Error loading .env.json file ", err)
	}
	if len(os.Args) < 3 {
		flag.Usage()
		os.Exit(1)
	}
	flag.Parse()
	if len(*symbol) < 1 {
		log.Fatal("No symbol specified. Exiting.")
	}
	bars := getBars(*symbol, *years)
	if len(bars) > 2 {
		fmt.Printf("%v period ADR : %.2f\n", 20, computeADR(bars[len(bars)-20:]))
		ranges := getRanges(bars)
		sort.Slice(ranges, func(i, j int) bool { return ranges[i].move < ranges[j].move })
		fmt.Printf("Largest gain on %v : %.2f (%.2f%%)\nLargest loss on %v : %.2f (%.2f%%)\n",
			ranges[len(ranges)-1].date.Format("2006-01-02"), ranges[len(ranges)-1].move, ranges[len(ranges)-1].percent,
			ranges[0].date.Format("2006-01-02"), ranges[0].move, ranges[0].percent)
	} else {
		fmt.Printf("Not enough bars returned: %v\n", len(bars))
	}
}
