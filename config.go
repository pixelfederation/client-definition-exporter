package main

// Ops represents the commandline/environment options for the program
type ConfigOpts struct {
    BaseDir         string `long:"basedir" short:"d" env:"BASEDIR" default:"/srv/www/" description:"Definition location"`
    MetricsBindAddr string `long:"metrics-bind-address" short:"b" env:"METRICS_BIND_ADDRESS" default:"9115" description:"Address for binding metrics listener"`
    MetricsPath     string `long:"metrics-path" env:"METRICS_PATH" default:"/metrics" description:"Metrics path"`
    Refresh         int `long:"refresh" short:"r" env:"REFRESH" default:"60" description:"Definition refresh rate"`
}