package main

import "time"

// Config and report logs files
const (
	StoreDir      = "stores/"
	ConfigFile    = StoreDir + "config.json"
	ReportLogFile = StoreDir + "reports.json"
)

// DataSetsCheckInterval dataset outliers checker interval
const DataSetsCheckInterval = 5 * time.Minute

// DateTimeFormat default date format
const DateTimeFormat = "2006-01-02 15:04:05"

// HTTP server params
const (
	ServerPort       = ":8086"
	HTTPReadTimeout  = 180 * time.Second
	HTTPWriteTimeout = 180 * time.Second
	MaxHearedBytes   = 1 << 20
)

// Outliers detection methods
const (
	ThreeSigmas = "3-sigmas"
)
