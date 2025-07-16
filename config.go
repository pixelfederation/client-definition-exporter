package main

// Ops represents the commandline/environment options for the program
type ConfigOpts struct {
    BaseDir         string `long:"basedir" short:"d" env:"BASEDIR" default:"/srv/www/" description:"Definition location"`
    MetricsBindAddr string `long:"metrics-bind-address" short:"b" env:"METRICS_BIND_ADDRESS" default:"9115" description:"Address for binding metrics listener"`
    MetricsPath     string `long:"metrics-path" short:"m" env:"METRICS_PATH" default:"/metrics" description:"Metrics path"`
    FilenameFormat  string `long:"filename-format" short:"f" env:"FILENAME_FORMAT" default:"client-data" description:"Definition files name format"`
    FileExtension   string `long:"file-extension" short:"e" env:"FILE_EXTENSION" default:"sqlite" description:"Definition files extension"`
    Refresh         int    `long:"refresh" short:"r" env:"REFRESH" default:"60" description:"Definition refresh rate"`
    LogLevel        int    `long:"loglevel" short:"v" env:"LOG_LEVEL" description:"log verbosity level for klog" default:"2"`
}
