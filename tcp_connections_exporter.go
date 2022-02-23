package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	tcp_connection_timewait = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connection_timewait",
		Help: "tcp_connection timewait count",
	})
	tcp_connection_established = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connection_established",
		Help: "tcp_connection established count",
	})
	tcp_connection_listen = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connection_listen",
		Help: "tcp_connection listen count",
	})
)

func init() {
	prometheus.MustRegister(tcp_connection_timewait)
	prometheus.MustRegister(tcp_connection_established)
	prometheus.MustRegister(tcp_connection_listen)
}

func connectionStatus(metric string) int {
	cmd := exec.Command("/bin/netstat", "-ant")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("%s", err)

	}
	result := strings.Count(string(out), metric)
	return result
}
func recordMetrics() {
	go func() {
		for {
			timewait := connectionStatus("TIME_WAIT")
			established := connectionStatus("ESTABLISHED")
			listen := connectionStatus("LISTEN")
			log.Info("Time Wait connections: " + strconv.Itoa(timewait))
			log.Info("Established connections: " + strconv.Itoa(established))
			log.Info("Listen connections: " + strconv.Itoa(listen))
			tcp_connection_established.Set(float64(established))
			tcp_connection_timewait.Set(float64(timewait))
			tcp_connection_timewait.Set(float64(listen))
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	bind := flag.String("bind", "0.0.0.0", "bind port default 0.0.0.0")
	port := flag.Int("port", 9319, "port to listen default 9319")
	flag.Parse()
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port " + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(*bind+":"+strconv.Itoa(*port), nil))
}
