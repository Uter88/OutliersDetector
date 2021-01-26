package main

import "time"

// JSONResponse default response
type JSONResponse struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// DataSetValues struct for DataSet values
type DataSetValues []DataSetValue

// OutliersDetection params
type OutliersDetection struct {
	OutliersMultipler       float64 `json:"OutliersMultipler"`
	StrongOutliersMultipler float64 `json:"StrongOutliersMultipler"`
}

// DataSetValue value for dataset
type DataSetValue struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// MetricValues struct for metric values
type MetricValues struct {
	Metric    string        `json:"Metric"`
	Attribute string        `json:"Arrtibute"`
	Values    DataSetValues `json:"values"`
}

// DataSet icoming data
type DataSet struct {
	SiteID                  string   `json:"siteId"`
	TimeAgo                 string   `json:"TimeAgo"`
	TimeStep                string   `json:"TimeStep"`
	OutliersDetectionMethod []string `json:"OutliersDetectionMethod"`
	MetricesList            []string `json:"MetricesList"`
	MinVisitorsPerTimeStep  int      `json:"MinVisitorsPerTimeStep"`
	OutliersDetection       `json:"OutliersDetection"`
	Metrics                 []MetricValues `json:"Values"`
}

// OutlierDetectResultRecord struct for outliers warnings and alarms detects
type OutlierDetectResultRecord struct {
	OutlierPeriodStart string `json:"OutlierPeriodStart"`
	OutlierPeriodEnd   string `json:"OutlierPeriodEnd"`
	Metric             string `json:"Metric"`
	Attribute          string `json:"Attribute"`
}

// OutliersDetectResult container for outliers warnings and alarms detects
type OutliersDetectResult struct {
	Warnings []OutlierDetectResultRecord `json:"Warnings"`
	Alarms   []OutlierDetectResultRecord `json:"Alarms"`
}

// OutlierDetectOutput output for DataSet outliers detection
type OutlierDetectOutput struct {
	SiteID                  string               `json:"siteId"`
	OutliersDetectionMethod string               `json:"OutliersDetectionMethod"`
	CheckTimeStart          string               `json:"checkTimeStart"`
	CheckTimeEnd            string               `json:"checkTimeEnd"`
	TimeAgo                 string               `json:"TimeAgo"`
	TimeStep                string               `json:"TimeStep"`
	DateStart               string               `json:"DateStart"`
	DateEnd                 string               `json:"DateEnd"`
	Result                  OutliersDetectResult `json:"Result"`
}

// OutliersResultLog outliers results logging
type OutliersResultLog struct {
	SiteID                  string `json:"siteId"`
	OutliersDetectionMethod string `json:"OutliersDetectionMethod"`
	TimeAgo                 string `json:"TimeAgo"`
	TimeStep                string `json:"TimeStep"`
	OutlierPeriodStart      string `json:"OutlierPeriodStart"`
	OutlierPeriodEnd        string `json:"OutlierPeriodEnd"`
	Metric                  string `json:"Metric"`
	Attribute               string `json:"Attribute"`
	Level                   string `json:"Level"`
}
