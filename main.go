package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

var (
	svStatRegEx  = regexp.MustCompile(`(up|down) \(.*\) (\d+) seconds(?:.*want (up|down))?`) // \([^\)]*\) (\d) seconds`)
	stateToFloat = map[string]float64{
		"up":   1.0,
		"down": 0.0,
	}
)

type exporter struct {
	serviceDir        string
	svStatBin         string
	metricUp          *prometheus.Desc
	metricWant        *prometheus.Desc
	metricStateChange *prometheus.Desc

	errors prometheus.Counter
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metricUp
	ch <- e.metricWant
	ch <- e.metricStateChange
	ch <- e.errors.Desc()
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	files, err := ioutil.ReadDir(e.serviceDir)
	if err != nil {
		log.Println("Couldn't read service directory:", err)
		e.errors.Inc()
		ch <- e.errors
		return
	}
	for _, file := range files {
		service := file.Name()
		if service == ".s6-svscan" || !file.IsDir() {
			continue
		}
		up, want, sc, err := e.svStat(service)
		if err != nil {
			log.Println(err)
			e.errors.Inc()
			continue
		}
		sc = float64(time.Now().Add(-time.Duration(sc) * time.Second).Unix())
		ch <- prometheus.MustNewConstMetric(e.metricUp, prometheus.GaugeValue, up, service)
		ch <- prometheus.MustNewConstMetric(e.metricWant, prometheus.GaugeValue, want, service)
		ch <- prometheus.MustNewConstMetric(e.metricStateChange, prometheus.GaugeValue, sc, service)
	}
	ch <- e.errors
}

func (e *exporter) svStat(name string) (up, want, sc float64, err error) {
	resp, err := exec.Command(e.svStatBin, filepath.Join(e.serviceDir, name)).Output()
	if err != nil {
		return up, want, sc, err
	}
	return parseSvStat(strings.TrimRight(string(resp), "\n\r"))
}

func parseSvStat(str string) (float64, float64, float64, error) {
	parts := svStatRegEx.FindStringSubmatch(str)
	if parts == nil {
		return 0.0, 0.0, 0.0, fmt.Errorf("Couldn't parse svstat response: '%s'", str)
	}
	up, ok := stateToFloat[parts[1]]
	if !ok {
		return 0.0, 0.0, 0.0, fmt.Errorf("Unknown state %s", parts[1])
	}
	sc, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	want := up
	if parts[3] != "" {
		want, ok = stateToFloat[parts[3]]
		if !ok {
			return 0.0, 0.0, 0.0, fmt.Errorf("Unknown state %s", parts[3])
		}
	}
	return up, want, float64(sc), nil
}

func main() {
	var (
		listenHTTP  = flag.String("h", ":9164", "Address to expose prometheus metrics on")
		serviceDir  = flag.String("d", "/etc/service", "Path to service directory")
		svStat      = flag.String("s", "s6-svstat", "svstat binary name")
		showVersion = flag.Bool("version", false, "Print version information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("s6_exporter"))
		os.Exit(0)
	}

	svStatBin, err := exec.LookPath(*svStat)
	if err != nil {
		log.Fatal(err)
	}

	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6_exporter_errors_total",
		Help: "Total number of errors s6_exporter encountered",
	})
	prometheus.MustRegister(version.NewCollector("s6_exporter"))
	prometheus.MustRegister(&exporter{
		serviceDir: *serviceDir,
		svStatBin:  svStatBin,
		errors:     errors,
		metricUp: prometheus.NewDesc(
			"s6_service_up",
			"State of s6 service, 1 = up, 0 = down",
			[]string{"service"}, nil,
		),
		metricWant: prometheus.NewDesc(
			"s6_service_wanted_up",
			"Wanted state of s6 service, 1 = up, 0 = down",
			[]string{"service"}, nil,
		),
		metricStateChange: prometheus.NewDesc(
			"s6_service_state_change_timestamp_seconds",
			"Unix timestamp of service's last state change.",
			[]string{"service"}, nil,
		),
	})

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html><head>s6_exporter</title></head><body>See /metrics</body></html>`))
		if err != nil {
			log.Println(err)
			errors.Inc()
		}
	})
	log.Println("Starting s6_exporter", version.Info())
	log.Println("Build context", version.BuildContext())
	log.Printf("Exposing metrics for %s on %s/metrics", *serviceDir, *listenHTTP)
	log.Fatal(http.ListenAndServe(*listenHTTP, nil))
}
