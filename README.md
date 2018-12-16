# loki_exporter

The [loki_exporter](https://github.com/ricoberger/loki_exporter) is a [Prometheus](https://prometheus.io) exporter for Loki. [Loki](https://github.com/grafana/loki) is a horizontally-scalable, highly-available, multi-tenant log aggregation system from the creators of [Grafana](https://grafana.com). The loki_exporter runs queries against the Loki API and returns the number of entries for each stream. This exporter is designed to detect critical log events, where the results can be used to create alerts in Prometheus.

## Building and running

### Local Build

```sh
make build
./bin/loki_exporter <flags>
```

Visiting http://localhost:9524/metrics will return metrics for all queries. The `loki_success` metric indicates if the execution of the queries succeeded.

### Building with Docker

```sh
docker build -t loki_exporter .
docker run -d -p 9524:9524 --name loki_exporter -v `pwd`:/config loki_exporter --config.file=/config/config.yml
```

## Usage and configuration

The loki_exporter is configured via a configuration file and command-line flags.

```
Usage of ./bin/loki_exporter:
  -config.file string
    	Configuration file in YAML format. (default "config.yml")
  -version
    	Show version information.
  -web.listen-address string
    	Address to listen on for web interface and telemetry. (default ":9524")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```

The configuration file is written in YAML format, defined by the scheme described below.

```yaml
# ======================== loki_exporter configuration =========================
#
# The loki_exporter is a Prometheus exporter for Loki. Loki is a
# horizontally-scalable, highly-available, multi-tenant log aggregation system
# from the creators of Grafana. The loki_exporter runs queries against the Loki
# API and returns the number of entries for each stream. This exporter is
# designed to detect critical log events, where the results can be used to
# create alerts in Prometheus.
#
# ------------------------------------ Loki ------------------------------------
#
loki:
  listenAddress: <string>
  basicAuth:
    enabled: <boolean>
    username: <string>
    password: <string>
#
# ---------------------------------- Queries -----------------------------------
#
queries:
  - name: <string>
    query: <string>
    limit: <integer>
    start: <string>
    end: <string>
    direction: <string>
    regexp: <string>
```

The configuration file is divided into two sections. The `loki` section is used for the configuration of the Loki API endpoint. The default value for the `listenAddress` is `http://localhost:3100`. Basic authentication is disabled by default.

The `queries` section represents all queries which should be run against the Loki API. The parameters and default values can be found in the following table:

| Parameter | Description | Default Value |
| --------- | ----------- | ------------- |
| `name` | A custom name for the query. The name is used for the exported metric: `loki_name` | |
| `query` | Query must be a logQL query. | |
| `limit` | Maximum number of entries which should be returned by the Loki API. | `-1` |
| `start` | The start time for the query. Must be a valid golang duration string. The duration is added to the current time. | `-24h` |
| `end` | The end time for the query. Must be a valid golang duration string. The duration is added to the current time. | `0s` |
| `direction` | Search direction must be `forward` or `backward`, useful when specifying a limit. | |
| `regexp` | A regular expression to filter the returned results. | |

## Example

The example is based on the official example from the [Loki documentation](https://github.com/grafana/loki). Clone the repository, go to the `examples` folder and run the `docker-compose` file. The `docker-compose` file will start a docker image for Loki, Promtail, Grafana, Prometheus and the loki_exporter.

```sh
git clone https://github.com/ricoberger/loki_exporter.git
cd loki_exporter/example
docker-compose up
```

Visiting the loki_exporter output on [http://localhost:9524](http://localhost:9524). The output will look like the one at the end of this section. You can also open the Prometheus dashboard on [http://localhost:9090](http://localhost:9090). There you can find the `loki_` metrics if in one of you log files under `/var/logs` an error occured.

```
# HELP loki_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which loki_exporter was built.
# TYPE loki_exporter_build_info gauge
loki_exporter_build_info{branch="HEAD",goversion="go1.11",revision="HEAD",version=""} 1
# HELP loki_exporter_total_scrapes Current total loki scrapes.
# TYPE loki_exporter_total_scrapes counter
loki_exporter_total_scrapes 3
# HELP loki_success Was the last scrape of loki successful.
# TYPE loki_success gauge
loki_success 1
# HELP loki_varlogs number of entries
# TYPE loki_varlogs gauge
loki_varlogs{filename="/var/log/docker.log",job="varlogs"} 151
loki_varlogs{filename="/var/log/docker.log.0",job="varlogs"} 33
loki_varlogs{filename="/var/log/kmsg.log",job="varlogs"} 2
loki_varlogs{filename="/var/log/vpnkit-forwarder.log",job="varlogs"} 175
loki_varlogs{filename="/var/log/vsudd.log",job="varlogs"} 6
loki_varlogs{filename="/var/log/vsudd.log.0",job="varlogs"} 6
```

## Dependencies

- [yaml.v2 - YAML support for the Go language](gopkg.in/yaml.v2)
- [Prometheus Go client library](github.com/prometheus/client_golang)
- [Common - Go libraries shared across Prometheus components and libraries](github.com/prometheus/common)
