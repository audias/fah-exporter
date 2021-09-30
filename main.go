package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultFahAddress = "192.168.0.201:36330"
)

var (
	fahAddress  = defaultFahAddress
	getAPI      = false
	apiThrottle time.Duration
	myClient    = &http.Client{Timeout: 10 * time.Second}
)

func main() {
	var (
		level              string
		listenAddress      string
		metricsPath        string
		socketActivate     bool
		noTimestamps       bool
		defaultThrottle, _ = time.ParseDuration("1h")
	)

	level = getVariable("LOG_LEVEL", "info")
	noTimestamps = getBooleanVariable("LOG_NO_TIMESTAMPS", false)
	listenAddress = getVariable("LISTEN_ADDRESS", "0.0.0.0:9659")
	fahAddress = getVariable("FAH_ADDRESS", defaultFahAddress)
	getAPI = getBooleanVariable("GET_DONOR_DATA", false)

	metricsPath = "/metrics"
	socketActivate = false
	getAPI = false
	apiThrottle = defaultThrottle

	setLogLevel(level)

	if noTimestamps || socketActivate {
		log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	} else {
		log.SetFormatter(&log.TextFormatter{DisableTimestamp: false, FullTimestamp: true})
	}

	prometheus.MustRegister(NewExporter())

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>FAH Exporter</title></head>
             <body>
             <h1>FAH Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
	     <h2>More information:</h2>
	     <p><a href="https://github.com/cosandr/fah-exporter">github.com/cosandr/fah-exporter</a></p>
             </body>
             </html>`))
	})

	listener := getListener(socketActivate, listenAddress)

	log.Infof("FAH client address: %s", fahAddress)
	log.Infof("Starting HTTP server on %s", listener.Addr().String())
	log.Fatal(http.Serve(listener, nil))
}

func setLogLevel(level string) {
	switch level {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln("Unrecognized minimum log level; using 'info' as default")
	}
}
