package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/hello-exporter/api"
)

const (
	metricsPrefix = "hello_"

	staleDataThreshold = 30 * time.Minute

	errorTypeOther = "other"
	errorTypeAuth  = "auth"
)

var (
	pillBatteryDesc = prometheus.NewDesc(
		metricsPrefix+"pill_battery_charge_percent",
		"Contains the battery charge of the sleeping pill.",
		[]string{"id", "color"}, nil)
	temperatureDesc = prometheus.NewDesc(
		metricsPrefix+"room_temperature_celsius",
		"Room temperature in degrees celsius.", nil, nil)
	noiseDesc = prometheus.NewDesc(
		metricsPrefix+"room_noise_decibel",
		"Room noise measurement in decibel.", nil, nil)
	humidityDesc = prometheus.NewDesc(
		metricsPrefix+"room_humidity_percent",
		"Room relative humidity in percent.", nil, nil)
	lightDesc = prometheus.NewDesc(
		metricsPrefix+"room_light_lux",
		"Room light level in lux.", nil, nil)
)

type helloCollector struct {
	client *api.HelloClient
	up     prometheus.Gauge
	time   prometheus.Histogram
	errors *prometheus.CounterVec
}

func newCollector(client *api.HelloClient) *helloCollector {
	return &helloCollector{
		client: client,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metricsPrefix + "up",
			Help: "Zero if there was an error during scrape process.",
		}),
		time: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: metricsPrefix + "scrape_duration_seconds",
			Help: "Contains the duration it took to scrape the Hello API.",
		}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: metricsPrefix + "errors_total",
			Help: "Counts the number of errors by type.",
		}, []string{"type"}),
	}
}

func (c *helloCollector) Describe(dChan chan<- *prometheus.Desc) {
	c.up.Describe(dChan)
	c.time.Describe(dChan)
	c.errors.Describe(dChan)
}

type collectFunc func(chan<- prometheus.Metric) error

func (c *helloCollector) Collect(mChan chan<- prometheus.Metric) {
	start := time.Now()
	up := true
	for _, collectFunc := range []collectFunc{
		c.collectDevices,
		c.collectRoomInfo,
	} {
		err := collectFunc(mChan)
		switch {
		case err == api.ErrWrongCredentials:
			up = false
			c.errors.WithLabelValues(errorTypeAuth).Inc()
			continue
		case err != nil:
			up = false
			c.errors.WithLabelValues(errorTypeOther).Inc()
			continue
		default:
		}
	}

	if up {
		c.up.Set(1)
	} else {
		c.up.Set(0)
	}

	c.time.Observe(time.Since(start).Seconds())

	c.up.Collect(mChan)
	c.time.Collect(mChan)
	c.errors.Collect(mChan)
}

func (c *helloCollector) collectDevices(mChan chan<- prometheus.Metric) error {
	devices, err := c.client.Devices()
	if err != nil {
		return err
	}

	collectPills(mChan, devices.Pills)
	return nil
}

func collectPills(mChan chan<- prometheus.Metric, pills []api.Pill) {
	for _, pill := range pills {
		sendMetric(mChan, pillBatteryDesc, prometheus.GaugeValue, float64(pill.BatteryLevel), pill.ID, pill.Color)
	}
}

func (c *helloCollector) collectRoomInfo(mChan chan<- prometheus.Metric) error {
	room, err := c.client.RoomInfo()
	if err != nil {
		return err
	}

	sendMetric(mChan, temperatureDesc, prometheus.GaugeValue, room.Temperature.Value)
	sendMetric(mChan, humidityDesc, prometheus.GaugeValue, room.Humidity.Value)
	sendMetric(mChan, noiseDesc, prometheus.GaugeValue, room.Sound.Value)
	sendMetric(mChan, lightDesc, prometheus.GaugeValue, room.Light.Value)
	return nil
}

func sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labels ...string) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, labels...)
	if err != nil {
		log.Printf("Error creating %s metric: %s", desc.String(), err)
	}
	ch <- m
}
