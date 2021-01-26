package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// GetDataSetBySiteID get single DataSet by siteID
func GetDataSetBySiteID(siteID string) (*DataSet, error) {
	sets, err := GetDataSets()

	if err != nil {
		return nil, err
	}
	for _, ds := range sets {
		if ds.SiteID == siteID {
			return &ds, nil
		}
	}
	return nil, errors.New("DataSet not found")
}

// GetDataSets get DataSets from store
func GetDataSets() ([]DataSet, error) {
	body, err := ReadFile(ConfigFile)

	if err != nil {
		return nil, err
	}

	dest := make(map[string][]DataSet)
	err = json.Unmarshal(body, &dest)

	if err != nil {
		return nil, fmt.Errorf("Error decode config file: %s", err)
	}
	if ds, ok := dest["Datasets"]; ok {
		return ds, nil
	}
	return nil, errors.New("Cannot found 'Datasets' key in config file root")
}

// ParseDuration parse time duration from string, like 1d, 24h
func ParseDuration(step string) (dur time.Duration, err error) {
	l := len(step)

	if l > 1 {
		val, err := strconv.Atoi(step[:l-1])

		if err != nil {
			return dur, fmt.Errorf("Error parse duration value: %s", err)
		}

		switch step[l-1] {
		case 'm':
			return time.Duration(val) * time.Minute, nil
		case 'h':
			return time.Duration(val) * time.Hour, nil
		case 'd':
			return time.Duration(val) * 24 * time.Hour, nil
		default:
			return dur, errors.New("Invalid duration, expected: m, h, d, w")
		}
	}
	return dur, errors.New("Corrupted duration param")
}

// OutliersReporter listens to the outliers channel, checks for uniqueness, in case of a new outliers - send a report
func OutliersReporter(c chan OutlierDetectOutput) {
	for o := range c {
		go LogOutliersReports(o)
	}
}

// GetReportOutliersLogs get outliers logs from file store
func GetReportOutliersLogs(o OutlierDetectOutput) (logs []OutliersResultLog, err error) {
	body, err := ReadFile(ReportLogFile)

	if err != nil {
		log.Printf("Error open report log file: %s\n", err.Error())
		return
	}
	dest := make(map[string][]OutliersResultLog)

	if err = json.Unmarshal(body, &dest); err != nil {
		log.Printf("Error decode reports log file: %s\n", err.Error())
		return
	}

	if err != nil {
		return
	}
	for _, l := range dest["Logs"] {
		if l.SiteID == o.SiteID && l.OutliersDetectionMethod == o.OutliersDetectionMethod {
			logs = append(logs, l)
		}
	}
	return
}

// CheckLogExists checks the result of determining new outliers for uniqueness
func CheckLogExists(logs []OutliersResultLog, rec OutlierDetectResultRecord) bool {
	for _, l := range logs {
		if rec.Equal(l) {
			return true
		}
	}
	return false
}

// LogOutliersReports get logs from store and check new outliers detection for unique
func LogOutliersReports(o OutlierDetectOutput) {
	logs, err := GetReportOutliersLogs(o)

	if err != nil {
		log.Printf("Error get logs from store: %s\n", err.Error())
		return
	}

	levels := map[string][]OutlierDetectResultRecord{
		"alarm":   o.Result.Alarms,
		"warning": o.Result.Warnings,
	}

	for lvl, recs := range levels {
		for _, rec := range recs {
			if !CheckLogExists(logs, rec) {
				WriteAndReportOutlierLog(o, rec, lvl)
			}
		}
	}
}

// WriteAndReportOutlierLog create outliers detection log, write to store and send report
func WriteAndReportOutlierLog(o OutlierDetectOutput, r OutlierDetectResultRecord, level string) {
	l := OutliersResultLog{
		SiteID:                  o.SiteID,
		OutliersDetectionMethod: o.OutliersDetectionMethod,
		TimeAgo:                 o.TimeAgo,
		TimeStep:                o.TimeStep,
		OutlierPeriodStart:      r.OutlierPeriodStart,
		OutlierPeriodEnd:        r.OutlierPeriodEnd,
		Metric:                  r.Metric,
		Attribute:               r.Attribute,
		Level:                   level,
	}
	if err := l.Save(); err != nil {
		log.Printf("Error save outliers log: %s\n", err.Error())
	}
	l.SendReport()
}

// DataSetsChecker check DataSets for outliers in a specific interval
func DataSetsChecker(c chan OutlierDetectOutput) {
	tiker := time.NewTicker(500 * time.Millisecond)
	defer tiker.Stop()

	for range tiker.C {
		datasets, err := GetDataSets()

		if err == nil {
			for _, ds := range datasets {
				ds.GenerateData()

				for _, o := range ds.DetectOutliers() {
					c <- o
				}
			}
		} else {
			log.Printf("Error get datasets: %s\n", err.Error())
		}
		time.Sleep(DataSetsCheckInterval)
	}
}
