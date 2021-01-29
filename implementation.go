package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// GetMeanStDev get mean and standart deviation values
func (dsv DataSetValues) GetMeanStDev() (float64, float64) {
	vals := make([]float64, dsv.Len())

	for i := range dsv {
		vals[i] = dsv[i].Value
	}
	return MeanStDev(vals...)
}

// GenerateData generate DataSet values
func (ds *DataSet) GenerateData() {
	end := time.Now()
	start := end.Add(-time.Hour * 24 * 35)

	for _, metric := range ds.MetricesList {
		mv := MetricValues{Metric: metric}

		for i := start.Unix(); i < end.Unix(); i += 60 * rand.Int63n(30) {
			dt := time.Unix(i, 0).UTC()
			val := GenerateValue(dt)
			mv.Values = append(mv.Values, DataSetValue{dt, val})
		}
		ds.Metrics = append(ds.Metrics, mv)
	}
}

// BreakIntoPieces break DataSetValues into pices by timeStep duration
func (dsv DataSetValues) BreakIntoPieces(timeStep time.Duration) (parts []DataSetValues) {
	var total = dsv.Len()

	var i, j int

	for i = 0; i < total; i++ {
		j = i
		next := dsv[i].Date.Add(timeStep).Round(timeStep)

		for j < total && dsv[j].Date.Before(next) {
			j++
		}
		parts = append(parts, dsv[i:j])
		i = j - 1
	}
	return
}

// Len return DataSetValues length
func (dsv DataSetValues) Len() int {
	return len(dsv)
}

// FilterByDatesAndMinValue filter dataset by start, end dates and minimal detect value
func (dsv DataSetValues) FilterByDatesAndMinValue(startDate, endDate time.Time, minValue float64) (values DataSetValues) {
	for _, v := range dsv {
		if v.Value >= minValue && v.Date.Unix() >= startDate.Unix() && v.Date.Unix() <= endDate.Unix() {
			values = append(values, v)
		}
	}
	return
}

// GetTimeAgoAndTimeStepDurations get TimeAgo and TimeStep durations
func (ds DataSet) GetTimeAgoAndTimeStepDurations() (timeAgo, timeStep time.Duration, err error) {
	timeAgo, err = ParseDuration(ds.TimeAgo)

	if err == nil {
		timeStep, err = ParseDuration(ds.TimeStep)
	}
	return
}

// MakeOutlierOutput returns OutlierDetectOutput
func (ds DataSet) MakeOutlierOutput(startDate, endDate time.Time) *OutlierDetectOutput {
	return &OutlierDetectOutput{
		SiteID:                  ds.SiteID,
		TimeAgo:                 ds.TimeAgo,
		TimeStep:                ds.TimeStep,
		OutliersDetectionMethod: ThreeSigmas,
		DateStart:               startDate.Format(DateTimeFormat),
		DateEnd:                 endDate.Format(DateTimeFormat),
		Result: OutliersDetectResult{
			Warnings: make([]OutlierDetectResultRecord, 0),
			Alarms:   make([]OutlierDetectResultRecord, 0),
		},
	}

}

// DetectOutliers detect DataSet values outliers
func (ds DataSet) DetectOutliers() (output []OutlierDetectOutput) {
	for _, method := range ds.OutliersDetectionMethod {
		switch method {

		case ThreeSigmas:
			for _, m := range ds.Metrics {
				out, err := ThreeSigmasOutlierDetector(ds, m)

				if err == nil {
					output = append(output, *out)
				}
			}
		default:
			log.Printf("Unsupported outlier detection method: %s\n", method)
		}
	}
	return
}

// SendReport send new outliers detection report
func (ol OutliersResultLog) SendReport() {
	// Do stuff
	msg := fmt.Sprintf(`
		Outliers detection result
		Start date: %s;
		End date: %s;
		Site ID: %s;
		Time ago: %s;
		Time step: %s;
		Metric: %s;
		Attribute: %s;
		Level: %s;
		Method: %s;
	`, ol.OutlierPeriodStart, ol.OutlierPeriodEnd, ol.SiteID, ol.TimeAgo,
		ol.TimeStep, ol.Metric, ol.Attribute, ol.Level, ol.OutliersDetectionMethod,
	)
	fmt.Println(msg)
}

// Save save outliers log to file
func (ol OutliersResultLog) Save() error {
	body, err := ReadFile(ReportLogFile)

	if err != nil {
		return fmt.Errorf("Error load reports log file: %s", err)
	}

	logs := make(map[string][]OutliersResultLog)

	if err = json.Unmarshal(body, &logs); err != nil {
		return fmt.Errorf("Error decode outliers log: %s", err)
	}

	if _, ok := logs["Logs"]; ok {
		logs["Logs"] = append(logs["Logs"], ol)
	}

	body, err = json.MarshalIndent(logs, "", " ")

	if err != nil {
		return fmt.Errorf("Error encode outliers log: %s", err)
	}

	if err = ioutil.WriteFile(ReportLogFile, body, 0644); err != nil {
		return fmt.Errorf("Error write outliers log to file: %s", err)
	}
	return nil
}

// Equal check equality of outliers detect result and log
func (odr OutlierDetectResultRecord) Equal(orl OutliersResultLog) bool {
	if odr.Metric != orl.Metric {
		return false
	}
	if odr.Attribute != orl.Attribute {
		return false
	}
	dates, err := ParseDates(
		odr.OutlierPeriodStart, odr.OutlierPeriodEnd,
		orl.OutlierPeriodStart, orl.OutlierPeriodEnd,
	)

	if err != nil {
		return false
	}
	start1, end1, start2, end2 := dates[0], dates[1], dates[2], dates[3]

	if start1.Day() == start2.Day() && end1.Day() == end2.Day() {
		if start1.Hour() == start2.Hour() && end1.Hour() == end2.Hour() {
			return true
		}
	}
	return false
}
