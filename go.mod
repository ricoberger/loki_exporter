module github.com/ricoberger/loki_exporter

require (
	github.com/grafana/loki v1.3.0
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/common v0.7.0
	gopkg.in/yaml.v2 v2.2.2
)

replace golang.org/x/net v0.0.0-20190813000000-74dc4d7220e7 => golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7

go 1.13
