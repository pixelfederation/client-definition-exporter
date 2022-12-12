package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

var FilenameFormat, FileExtension = "client - data", "sqlite"

func visit(path string, di fs.DirEntry, err error) error {
	// r, _ := regexp.Compile("client-data.*.sqlite")
	r, _ := regexp.Compile(FilenameFormat + ".*." + FileExtension)

	if r.MatchString(path) {
		klog.Info("Visited: %s\n", filepath.Base(path))

		decomposed_f := strings.Split(path, "/")
		dynamic := decomposed_f[len(decomposed_f)-2]

		if dynamic == "client-resources" {
			dynamic = "none"
		}

		if decomposed_f[4] == "" {
			klog.Fatalf("Error decompose path: %v", err)
		}

		clienDefinitionMetrics(prometheus.Labels{
			"origin":  decomposed_f[4],
			"env":     decomposed_f[3],
			"dynamic": dynamic,
			// "version": strings.ReplaceAll(filepath.Base(path), ".sqlite", "")})
			"version": strings.ReplaceAll(filepath.Base(path), "."+FileExtension, "")})
	}
	return nil
}

func clienDefinitionMetrics(labels map[string]string) {
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   "client",
			Name:        "definition_data_version",
			Help:        "Client definition data version.",
			ConstLabels: labels,
		},
		func() float64 { return float64(1) },
	)); err == nil {
		klog.Info("GaugeFunc 'definition_data_version' registered with labels: %s \n", labels)
	}
}

func recordMetrics(conf *ConfigOpts) {
	ticker := time.NewTicker(time.Duration(conf.Refresh) * time.Second)
	for {
		err := filepath.WalkDir(conf.BaseDir, visit)
		fmt.Printf("filepath.WalkDir() returned %v\n", err)
		<-ticker.C
	}
}

func main() {

	klog.InitFlags(nil)
	conf := &ConfigOpts{}
	parser := flags.NewParser(conf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		klog.Fatalf("Error parsing flags: %v", err)
	}

	FileExtension = conf.FileExtension
	FilenameFormat = conf.FilenameFormat

	go recordMetrics(conf)

	http.Handle(conf.MetricsPath, promhttp.Handler())
	http.ListenAndServe(":"+conf.MetricsBindAddr, nil)
}
