package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type PrometheusCollector struct {
	Name                string
	HelpMsg             string
	Labels              prometheus.Labels
	Url                 string
	IsRespTimeCollector bool
	PcHttpClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (pc *PrometheusCollector) CollectFunc() float64 {
	start := time.Now()
	request, err := http.NewRequest(http.MethodGet, pc.Url, nil)
	if err != nil {
		return 0
	}
	response, err := pc.PcHttpClient.Do(request)
	if pc.IsRespTimeCollector {
		elapsed := time.Since(start).Milliseconds()
		return float64(elapsed)
	}

	if err == nil && response.StatusCode == http.StatusOK {
		return 1
	} else {
		return 0
	}
}

func registerCollector(aInCollector PrometheusCollector, aInRegistry *prometheus.Registry) error {
	if err := aInRegistry.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        aInCollector.Name,
			Help:        aInCollector.HelpMsg,
			ConstLabels: aInCollector.Labels,
		},
		aInCollector.CollectFunc,
	)); err == nil {
		fmt.Printf("GaugeFunc %s registered.\n", aInCollector.Name)
		return nil
	} else {
		return err
	}
}

func initCollectors() []PrometheusCollector {
	var httpClient HttpClient
	httpClient = &http.Client{}

	http200StatusCollector := PrometheusCollector{
		"sample_external_url_up_200",
		"Up/Down status of https://httpstat.us/200.",
		prometheus.Labels{"url": "https://httpstat.us/200"},
		"https://httpstat.us/200", false, httpClient}

	http200RespTimeCollector := PrometheusCollector{
		"sample_external_url_response_ms_200",
		"Response time of https://httpstat.us/200 in milliseconds.",
		prometheus.Labels{"url": "https://httpstat.us/200"},
		"https://httpstat.us/200", true, httpClient}

	http503StatusCollector := PrometheusCollector{
		"sample_external_url_up_503",
		"Up/Down status of https://httpstat.us/503.",
		prometheus.Labels{"url": "https://httpstat.us/503"},
		"https://httpstat.us/503", false, httpClient}

	http503RespTimeCollector := PrometheusCollector{
		"sample_external_url_response_ms_503",
		"Response time of https://httpstat.us/503 in milliseconds.",
		prometheus.Labels{"url": "https://httpstat.us/503"},
		"https://httpstat.us/503", true, httpClient}

	collectors := []PrometheusCollector{http200StatusCollector, http200RespTimeCollector, http503StatusCollector, http503RespTimeCollector}
	return collectors
}

func initMetrics(aInRegistry *prometheus.Registry) error {
	collectors := initCollectors()
	for _, collector := range collectors {
		err := registerCollector(collector, aInRegistry)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	r := prometheus.NewRegistry()
	err := initMetrics(r)
	if err != nil {
		fmt.Errorf("Failed to register prometheus metrics with error: %s", err.Error())
		return
	}
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	http.ListenAndServe(":2112", nil)
}
