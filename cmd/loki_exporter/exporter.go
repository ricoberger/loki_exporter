package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const namespace = "loki"

// Exporter represents the structur of the exporter
type Exporter struct {
	up           prometheus.Gauge
	totalScrapes prometheus.Counter
	lokiMetrics  map[string]*prometheus.GaugeVec
	client       *client.Client
}

// NewExporter returns an initialized exporter
func NewExporter(lokiMetrics map[string]*prometheus.GaugeVec) (*Exporter, error) {
	client := &client.Client{
		Address: exporterConfig.Loki.ListenAddress,
	}
	if exporterConfig.Loki.BasicAuth.Enabled {
		client.Username = exporterConfig.Loki.BasicAuth.Username
		client.Password = exporterConfig.Loki.BasicAuth.Password
	}
	return &Exporter{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "success",
			Help:      "Was the last scrape of loki successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total loki scrapes.",
		}),
		lokiMetrics: lokiMetrics,
		client:      client,
	}, nil
}

func (e *Exporter) scrape() {
	e.totalScrapes.Inc()
	e.up.Set(0)

	// Labels
	if exporterConfig.Metrics.Labels == true {
		labels, err := e.getLabels()
		if err != nil {
			log.Errorln(err)
		} else {
			// Label Values
			if exporterConfig.Metrics.LabelValues == true {
				err := e.getLabelValues(labels)
				if err != nil {
					log.Errorln(err)
				}
			}
		}
	}

	// Queries
	if exporterConfig.Metrics.Queries == true {
		err := e.getQueries()
		if err != nil {
			log.Errorln(err)
		}
	}

	e.up.Set(1)
}

func (e *Exporter) resetMetrics() {
	for _, m := range e.lokiMetrics {
		m.Reset()
	}
}

func (e *Exporter) collectMetrics(metrics chan<- prometheus.Metric) {
	for _, m := range e.lokiMetrics {
		m.Collect(metrics)
	}
}

// Describe describes all the metrics ever exported by the loki_exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
}

// Collect queries the Loki API and delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.resetMetrics()
	e.scrape()

	ch <- e.up
	ch <- e.totalScrapes
	e.collectMetrics(ch)
}

func (e *Exporter) getLabels() (*loghttp.LabelResponse, error) {
	res, err := e.client.ListLabelNames(true)
	if err != nil {
		return nil, err
	}

	e.lokiMetrics["labels_total"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "labels_total",
			Help:      "total number of label",
		},
		[]string{},
	)

	e.lokiMetrics["labels_total"].With(prometheus.Labels{}).Set(float64(len(res.Data)))

	return res, nil
}

func (e *Exporter) getLabelValues(labels *loghttp.LabelResponse) error {
	for _, label := range labels.Data {
		res, err := e.client.ListLabelValues(label, true)
		if err != nil {
			return err
		}

		e.lokiMetrics["label_values_"+label+"_total"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "label_values_total",
				Help:      "total number of label values",
			},
			[]string{
				"label",
			},
		)

		e.lokiMetrics["label_values_"+label+"_total"].With(prometheus.Labels{"label": label}).Set(float64(len(res.Data)))
	}

	return nil
}

func (e *Exporter) getQueries() error {
	for _, query := range exporterConfig.Queries {
		res, err := e.client.QueryRange(
			query.Query,
			query.Limit,
			getQueryTime(query.Start),
			getQueryTime(query.End),
			getDirection(query.Direction),
			// let loki choose a default step
			0,
			// quiet logging
			true,
		)
		if err != nil {
			return err
		}
		if res.Data.ResultType != logql.ValueTypeStreams {
			return errors.New("invalid result type")
		}

		for index, stream := range res.Data.Result.(loghttp.Streams) {
			name := query.Name + strconv.FormatInt(int64(index), 10)
			labelNames := make([]string, 0, len(stream.Labels.Map()))

			for k := range stream.Labels.Map() {
				labelNames = append(labelNames, k)
			}

			e.lokiMetrics[name] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "query_" + query.Name + "_total",
					Help:      "number of entries",
				},
				labelNames,
			)

			e.lokiMetrics[name].With(stream.Labels.Map()).Set(float64(len(stream.Entries)))
		}
	}

	return nil
}

func getQueryTime(t time.Duration) time.Time {
	return time.Now().Add(t)
}

func getDirection(direction string) logproto.Direction {
	if direction == "forward" {
		return logproto.FORWARD
	}
	return logproto.BACKWARD
}
