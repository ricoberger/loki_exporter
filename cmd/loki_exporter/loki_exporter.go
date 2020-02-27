package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ricoberger/loki_exporter/pkg/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	exporterConfig config.Config

	listenAddress = flag.String("web.listen-address", ":9524", "Address to listen on for web interface and telemetry.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	showVersion   = flag.Bool("version", false, "Show version information.")
	configFile    = flag.String("config.file", "config.yml", "Configuration file in YAML format.")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Load configuration file
	err := exporterConfig.LoadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Show version information
	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("loki_exporter"))
		os.Exit(0)
	}

	// Prometheus exporter
	log.Infoln("Starting loki_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	metrics := map[string]*prometheus.GaugeVec{}
	exporter, err := NewExporter(metrics)
	if err != nil {
		log.Fatal(err)
	}

	// use our own registry to not get metrics from loki dependencies like etcd
	r := prometheus.NewRegistry()
	r.MustRegister(exporter)
	r.MustRegister(version.NewCollector("loki_exporter"))
	r.MustRegister(prometheus.NewGoCollector())
	r.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	handler := promhttp.InstrumentMetricHandler(r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Loki Exporter</title></head>
		<body>
		<h1>Loki Exporter</h1>
		<p><a href='` + *metricsPath + `'>Metrics</a></p>
		<p><ul>
		<li>version: ` + version.Version + `</li>
		<li>branch: ` + version.Branch + `</li>
		<li>revision: ` + version.Revision + `</li>
		<li>go version: ` + version.GoVersion + `</li>
		<li>build user: ` + version.BuildUser + `</li>
		<li>build date: ` + version.BuildDate + `</li>
		</ul></p>
		</body>
		</html>`))
	})
	log.Fatalln(http.ListenAndServe(*listenAddress, nil))
}
