package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "jstat"
)

var (
	listenAddress = flag.String("web.listen-address", ":9010", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	jstatPath     = flag.String("jstat.path", "/usr/bin/jstat", "jstat path")
	targetPid     = flag.String("target.pid", ":0", "target pid")
)

type Exporter struct {
	jstatPath string
	targetPid string
	GaugeMap  map[string]prometheus.Gauge
}

func NewExporter(jstatPath string, targetPid string) *Exporter {
	return &Exporter{
		jstatPath: jstatPath,
		targetPid: targetPid,
		GaugeMap:  make(map[string]prometheus.Gauge),
	}
}

// jstat -gccapacity
var GccapacityFields = []string{"ngcmn", "ngcmx", "ngc", "s0c", "s1c", "ec", "ogcmn", "ogcmx", "ogc", "oc", "mcmn", "mcmx", "mc", "ccsmn", "ccsmx", "ccsc", "ygc", "fgc"}

// jstat -gc
// var GcFields = []string{"s0c", "s1c", "s0u", "s1u", "ec", "eu", "oc", "ou", "mc", "mu", "ccsc", "ccsu", "ygc", "ygct", "fgc", "fgct", "gct"}
var GcFields = []string{"", "", "s0u", "s1u", "", "eu", "", "ou", "", "mu", "", "ccsu", "", "ygct", "", "fgct", "gct"}

// jstat -gcold
// var GcoldFields = []string{"mc", "mu", "ccsc", "ccsu", "oc", "ou", "ygc", "fgc", "fgct", "gct"}
var GcoldFields = []string{}

// jstat -gcnew
// var GcnewFields = []string{"s0c", "s1c", "s0u", "s1u", "tt", "mtt", "dss", "ec", "eu", "ygc", "ygct"}
var GcnewFields = []string{"", "", "", "", "tt", "mtt", "dss", "", "", "", ""}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range e.GaugeMap {
		v.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.JstatGccapacity(ch)
	// e.JstatGcold(ch)
	e.JstatGcnew(ch)
	e.JstatGc(ch)
}

func (e *Exporter) JstatGccapacity(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gccapacity", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i != 1 {
			continue
		}
		parts := strings.Fields(line)
		for j, field := range GccapacityFields {
			if field == "" {
				continue
			}
			value, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				log.Fatal(err)
			}
			v := e.GaugeMap[field]
			if v == nil {
				v = prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      field,
					Help:      field,
				})
			}
			v.Set(value)
			v.Collect(ch)
		}
	}
}

func (e *Exporter) JstatGcold(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gcold", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i != 1 {
			continue
		}

		parts := strings.Fields(line)
		for j, field := range GcoldFields {
			if field == "" {
				continue
			}

			value, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				log.Fatal(err)
			}
			v := e.GaugeMap[field]
			if v == nil {
				v = prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      field,
					Help:      field,
				})
			}
			v.Set(value)
			v.Collect(ch)
		}
	}
}

func (e *Exporter) JstatGcnew(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gcnew", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i != 1 {
			continue
		}

		parts := strings.Fields(line)
		for j, field := range GcnewFields {
			if field == "" {
				continue
			}
			value, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				log.Fatal(err)
			}
			v := e.GaugeMap[field]
			if v == nil {
				v = prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      field,
					Help:      field,
				})
			}
			v.Set(value)
			v.Collect(ch)
		}
	}
}

func (e *Exporter) JstatGc(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gc", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i != 1 {
			continue
		}

		parts := strings.Fields(line)
		for j, field := range GcFields {
			if field == "" {
				continue
			}
			value, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				log.Fatal(err)
			}
			v := e.GaugeMap[field]
			if v == nil {
				v = prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      field,
					Help:      field,
				})
			}
			v.Set(value)
			v.Collect(ch)
		}
	}
}

func main() {
	flag.Parse()

	if _, err := exec.LookPath(*jstatPath); err != nil {
		log.Fatalf("jstat not found at %s", *jstatPath)
		return
	}

	exporter := NewExporter(*jstatPath, *targetPid)
	prometheus.MustRegister(exporter)
	prometheus.Unregister(collectors.NewGoCollector())

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>jstat Exporter</title></head>
		<body>
		<h1>jstat Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
