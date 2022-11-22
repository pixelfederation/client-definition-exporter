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
	"k8s.io/klog/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func visit(path string, di fs.DirEntry, err error) error {
	r, _ := regexp.Compile("client-data.*.sqlite")
	
	if r.MatchString(path) {
		fmt.Printf("Visited: %s\n", filepath.Base(path))
		
		decomposed_f := strings.Split(path, "/")
		dynamic := decomposed_f[len(decomposed_f)-2]

		if dynamic == "client-resources" {
		    dynamic = "none"
		}
		
		clienDefinitionMetrics(prometheus.Labels{
			"origin": decomposed_f[4],
			"env": decomposed_f[3],
			"dynamic": dynamic,
			"version": strings.ReplaceAll(filepath.Base(path), ".sqlite", "")})
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
		fmt.Printf("GaugeFunc 'definition_data_version' registered with labels: %s \n", labels)
	}
}

func recordMetrics(conf *ConfigOpts) {
	go func() {
		for {
			err := filepath.WalkDir(conf.BaseDir, visit)
			fmt.Printf("filepath.WalkDir() returned %v\n", err)
			time.Sleep(time.Duration(conf.Refresh) * time.Second)
		}
	}()
}

func main() {

	klog.InitFlags(nil)
	conf := &ConfigOpts{}
	parser := flags.NewParser(conf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		klog.Fatalf("Error parsing flags: %v", err)
	}

	recordMetrics(conf)

	http.Handle(conf.MetricsPath, promhttp.Handler())
	http.ListenAndServe(":" + conf.MetricsBindAddr, nil)
}