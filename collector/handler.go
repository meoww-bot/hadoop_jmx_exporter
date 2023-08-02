package collector

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(w http.ResponseWriter, r *http.Request, logger log.Logger) {

	params := r.URL.Query()

	exportSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hadoop_jmx_export_success",
		Help: "Displays whether or not the exporter was a success",
	})
	exportDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hadoop_jmx_export_duration_seconds",
		Help: "Returns how long the exporter took to complete in seconds",
	})

	targetParam := params.Get("target")
	if targetParam == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	target, err := url.QueryUnescape(targetParam)
	if err != nil {
		http.Error(w, "Error decoding target parameter", http.StatusBadRequest)
		level.Error(logger).Log("msg", "Error decoding target parameter", "err", err)

		// log.Error("Error decoding target parameter:", err)
		return
	}

	t := Target{
		Url:    target,
		Logger: logger,
	}

	KrbPrincipalParam := params.Get("principal")

	if KrbPrincipalParam != "" {
		KrbPrincipal, err := url.QueryUnescape(KrbPrincipalParam)
		if err != nil {
			http.Error(w, "Error decoding principal parameter", http.StatusBadRequest)
			level.Error(logger).Log("msg", "Error decoding principal parameter", "err", err)

			// log.Error("Error decoding principal parameter:", err)
		}

		if KrbPrincipal != "" {
			t.KrbPrincipal = KrbPrincipal
		}
	}

	KrbPasswordParam := params.Get("password")

	if KrbPasswordParam != "" {
		KrbPassword, err := url.QueryUnescape(KrbPasswordParam)
		if err != nil {
			http.Error(w, "Error decoding password parameter", http.StatusBadRequest)
			level.Error(logger).Log("msg", "Error decoding password parameter", "err", err)

			// log.Error("Error decoding password parameter:", err)
		}

		if KrbPassword != "" {
			t.KrbAuthMethod = "password"
			t.KrbPassword = KrbPassword
		}
	}

	KrbKtPathParam := params.Get("ktpath")

	if KrbKtPathParam != "" {
		KrbKtPath, err := url.QueryUnescape(KrbKtPathParam)

		if err != nil {
			http.Error(w, "Error decoding ktpath parameter", http.StatusBadRequest)
			level.Error(logger).Log("msg", "Error decoding ktpath parameter", "err", err)

			// log.Error("Error decoding ktpath parameter:", err)
		}
		if KrbKtPath != "" {
			t.KrbAuthMethod = "keytab"
			t.KrbKtPath = KrbKtPath
		}
	}

	err = t.getCollectorName()

	if err != nil {
		exportSuccessGauge.Set(0)
		level.Error(logger).Log("msg", "Error get collector name", "err", err)

	}

	exporter, ok := Collectors[t.ExporterName]
	if !ok {
		exportSuccessGauge.Set(0)
		// http.Error(w, fmt.Sprintf("Unknown exporter %q,  Http StatusText: %s", t.ExporterName, t.RespStatus), http.StatusBadRequest)

	}

	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(exportSuccessGauge)
	registry.MustRegister(exportDurationGauge)

	if ok {
		level.Info(logger).Log("target", target, "collector", t.ExporterName)

		success := exporter(t, registry)
		duration := time.Since(start).Seconds()
		exportDurationGauge.Set(duration)
		if success {
			exportSuccessGauge.Set(1)
		} else {
			exportSuccessGauge.Set(0)
		}

	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
