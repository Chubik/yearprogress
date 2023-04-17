package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
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
	cpem, kpem        string
)

const (
	certPath = "/home/gentle/cert/"
)

func main() {

	cpem, err := url.JoinPath(certPath, "cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	kpem, err := url.JoinPath(certPath, "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	useTLS := flag.Bool("tls", false, "Enable TLS")
	certFile := flag.String("cert", cpem, "Path to SSL certificate")
	keyFile := flag.String("key", kpem, "Path to SSL key")
	port := flag.Int("port", 8085, "Port to listen on")
	flag.Parse()

	currentYear = time.Now().Year()
	statistics = &StatData{}

	router := httprouter.New()
	router.GET("/", sayGen)
	router.GET("/len/:len", sayGen)
	router.GET("/stat", stat)

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	if *useTLS {
		config := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		server.TLSConfig = config

		log.Printf("Starting server with TLS on port %d...\n", *port)
		log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
	} else {
		log.Printf("Starting server without TLS on port %d...\n", *port)
		log.Fatal(server.ListenAndServe())
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
