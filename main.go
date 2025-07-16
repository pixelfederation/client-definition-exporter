package main

import (
    "flag"
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
        klog.V(4).Info("Visited: %s\n", filepath.Base(path))

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
        klog.V(4).Info("GaugeFunc 'definition_data_version' registered with labels: %s \n", labels)
    }
}

func recordMetrics(conf *ConfigOpts) {
    ticker := time.NewTicker(time.Duration(conf.Refresh) * time.Second)
    for {
        err := filepath.WalkDir(conf.BaseDir, visit)
        if err != nil {
            fmt.Printf("filepath.WalkDir() returned %v\n", err)
        }
        <-ticker.C
    }
}

func main() {

    klog.InitFlags(nil)

    var conf ConfigOpts
    _, err := flags.Parse(&conf)

    if err != nil {
      panic("no arguments")
    }

    flag.Set("v", fmt.Sprintf("%d", conf.LogLevel))

    defer klog.Flush()

    FileExtension = conf.FileExtension
    FilenameFormat = conf.FilenameFormat

    go recordMetrics(&conf)

    http.Handle(conf.MetricsPath, promhttp.Handler())
    err = http.ListenAndServe(":"+conf.MetricsBindAddr, nil)
    if err != nil {
        klog.Fatalf("Failed to start server: %v", err)
    }
}
