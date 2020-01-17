package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	service service
}

var validationSuccessfulGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "trivrost",
		Name:      "validation_ok",
		Help:      "1 if validation ok, 0 otherwise",
	},
)

var promHttpHandler = promhttp.Handler()

func registerMetrics() {
	if err := prometheus.Register(validationSuccessfulGauge); err != nil {
		panic(fmt.Errorf("registering validationSuccessfulGauge failed: %w", err))
	}
}

func (s metrics) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	configUrl := getConfigUrl(req, s.service.flags)
	if configUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Required query parameter '"+configUrlParameterName+"' is missing.\n")
		return
	}
	reps := validateDeploymentConfig(configUrl, s.service.flags.SkipUrlCheck, s.service.flags.SkipJarChek)

	if reps.HaveError() {
		logReports(reps, true)
		validationSuccessfulGauge.Set(0)
	} else {
		validationSuccessfulGauge.Set(1)
	}

	promHttpHandler.ServeHTTP(w, req)
}
