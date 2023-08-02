package main

import (
	"hadoop_jmx_exporter/collector"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	_ "github.com/sijms/go-ora/v2"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	// Required for debugging
	// _ "net/http/pprof"
)

var (
	// Version will be set at build time.
	Version      = "0.0.0.dev"
	scrapePath   = kingpin.Flag("web.scrape-path", "Path under which to expose metrics. (env: TELEMETRY_PATH)").Default(getEnv("TELEMETRY_PATH", "/scrape")).String()
	toolkitFlags = webflag.AddFlags(kingpin.CommandLine, ":9070")
)

func main() {
	promLogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promLogConfig)
	kingpin.HelpFlag.Short('\n')
	kingpin.Version(version.Print("hadoop_jmx_exporter"))
	kingpin.Parse()
	logger := promlog.New(promLogConfig)

	version.Version = Version

	level.Info(logger).Log("msg", "Starting hadoop_jmx_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build", version.BuildContext())

	http.HandleFunc("/scrape", scrapeHandle(logger))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><head><title>Hadoop Jmx Exporter " + Version + "</title></head><body><h1>Hadoop Jmx Exporter " + Version + "</h1><p><a href='" + *scrapePath + "'>Scrape</a></p></body></html>"))
	})

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Listening error", "reason", err)
		os.Exit(1)
	}
}

func scrapeHandle(logger log.Logger) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		collector.Handler(w, r, logger)

	}
}

// getEnv returns the value of an environment variable, or returns the provided fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
