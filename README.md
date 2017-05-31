# hello-exporter

Simple [prometheus](https://prometheus.io) exporter for getting data from [Hello Sense](https://hello.is) devices into prometheus.

## Installation

If you have a working Go installation, getting the binary should be as simple as

```bash
go get github.com/xperimental/hello-exporter
```

There is also a `build-arm.sh` script if you want to run the exporter on an ARMv7 device.

## Usage

```
$ hello-exporter --help
Usage of hello-exporter:
  -a, --addr string       Address to listen on. (default ":9258")
  -p, --password string   Password of Hello account.
  -u, --username string   Username of Hello account.
```

After starting the server will offer the metrics on the `/metrics` endpoint, which can be used as a target for prometheus.

The exporter will query the Hello API every time it is scraped by prometheus. It does not make sense to scrape the Hello API with a small interval as the sensors only update their data every few minutes, so don't forget to set a slower scrape interval for this exporter:

```yml
scrape_configs:
  - job_name: 'hello'
    scrape_interval: 90s
    static_configs:
      - targets: ['localhost:9258']
```
