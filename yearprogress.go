package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Percent int    `json:"percent"`
	Str     string `json:"str"`
}

var (
	ct     time.Time          //current time
	Start  time.Time          //start date of the year
	wport  string    = "8085" //default api port, should be changed by ENV apram PORT
	bg     string    = "▓"    //progress symbol
	pr     string    = "░"    //background symbol
	maxStr int       = 20     //max symbols for generated string
)

const (
	yDays = 365 //days in year
)

func main() {
	format := "2006-01-02 15:04:05"
	Start, _ = time.Parse(format, "2019-01-01 00:00:00")
	http.HandleFunc("/", sayGen)
	if err := http.ListenAndServe(GetPort(), nil); err != nil {
		panic(err)
	}

}

func genStr(p, len int) string {
	s := ""
	f := (p * len) / 100
	for i := 0; i < len; i++ {
		if f <= i {
			s += pr
		} else {
			s += bg
		}
	}
	return s
}

func sayGen(w http.ResponseWriter, r *http.Request) {
	rt := ReturnPercent(Start)
	ri := int(math.Round(rt))
	str := genStr(ri, maxStr)

	resp := Response{
		Percent: ri,
		Str:     str,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func PercentOf(current int, all int) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return percent
}

func ReturnPercent(start time.Time) float64 {
	ct = time.Now()
	diff := ct.Sub(start)
	dLeft := int(diff.Hours() / 24)
	perc := PercentOf(dLeft, yDays)
	return perc
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = wport
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
