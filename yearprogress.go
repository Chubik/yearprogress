package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Percent int    `json:"percent"`
	Str     string `json:"str"`
}

type StatData struct {
	Visitors int `json:"visitors"`
}

var (
	currentYear       int
	progressBarLength int    = 20
	bg                string = "▓"
	pr                string = "░"
	statistics        *StatData
)

func main() {
	currentYear = time.Now().Year()
	statistics = &StatData{}

	router := httprouter.New()
	router.GET("/", sayGen)
	router.GET("/len/:len", sayGen)
	router.GET("/stat", stat)

	if err := http.ListenAndServe(GetPort(), router); err != nil {
		panic(err)
	}
}

func genStr(p, length int) string {
	var builder strings.Builder
	filled := (p * length) / 100

	for i := 0; i < length; i++ {
		if i < filled {
			builder.WriteString(bg)
		} else {
			builder.WriteString(pr)
		}
	}
	return builder.String()
}

func stat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	st, err := json.Marshal(statistics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(st)
}

func sayGen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	statistics.Visitors++
	ls, err := strconv.Atoi(ps.ByName("len"))
	if err != nil {
		ls = progressBarLength
	}

	percent := yearProgressPercentage()
	str := genStr(percent, ls)

	resp := Response{
		Percent: percent,
		Str:     str,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(js)
}

func yearProgressPercentage() int {
	startOfYear := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	daysPassed := now.Sub(startOfYear).Hours() / 24
	return int(math.Round((daysPassed / 365) * 100))
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
