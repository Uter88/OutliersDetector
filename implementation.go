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
	vals := make([]float64, len(dsv))

	for i := range dsv {
		vals[i] = dsv[i].Value
	}
	return MeanStDev(vals...)
}

// GetDateRanges generate date ranges
func (ds DataSet) GetDateRanges() (start, end time.Time, err error) {
	dur, err := ParseDuration(ds.TimeAgo)

	if err != nil {
		return
	}
	end = time.Now().UTC()
	start = end.Add(-dur)
	return start, end, nil
}

// GenerateData make DataSet values
func (ds *DataSet) GenerateData() {
	start, end, err := ds.GetDateRanges()

	if err != nil {
		log.Printf("Error generate data: %s\n", err.Error())
		return
	}
	for _, metric := range ds.MetricesList {
		mv := MetricValues{Metric: metric}

		for i := start.Unix(); i < end.Unix(); i += 60 * rand.Int63n(30) {
			dt := time.Unix(i, 0)
			val := GenerateValue(dt)
			mv.Values = append(mv.Values, DataSetValue{dt, val})
		}
		ds.Metrics = append(ds.Metrics, mv)
	}
}

// FilterByMinValue filter dataset values by minimal value
func (dsv DataSetValues) FilterByMinValue(min float64) (vals DataSetValues) {
	for _, v := range dsv {
		if v.Value >= min {
			vals = append(vals, v)
		}
	}
	return
}

// GetMinValue get outliers detection minimal value
func (ds DataSet) GetMinValue() (v float64) {
	dur, err := ParseDuration(ds.TimeAgo)

	if err != nil {
		return float64(ds.MinVisitorsPerTimeStep)
	}
	step, err := ParseDuration(ds.TimeStep)

	if err != nil {
		return float64(ds.MinVisitorsPerTimeStep)
	}
	return float64(dur/step) * float64(ds.MinVisitorsPerTimeStep)
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
func (or OutlierDetectResultRecord) Equal(o OutliersResultLog) bool {
	if or.Metric != o.Metric {
		return false
	}
	if or.Attribute != o.Attribute {
		return false
	}
	start1, end1, err1 := or.ParseDates()
	start2, end2, err2 := o.ParseDates()

	if err1 != nil || err2 != nil {
		return false
	}
	if start1.Day() == start2.Day() && end1.Day() == end2.Day() {
		if start1.Hour() == start2.Hour() && end1.Hour() == end2.Hour() {
			return true
		}
	}
	return false
}

// ParseDates parse DateStart and DateEnd from string
func (ol OutliersResultLog) ParseDates() (start, end time.Time, err error) {
	start, err = time.Parse(DateTimeFormat, ol.OutlierPeriodStart)

	if err != nil {
		return
	}
	end, err = time.Parse(DateTimeFormat, ol.OutlierPeriodEnd)
	return
}

// ParseDates parse DateStart and DateEnd from string
func (or OutlierDetectResultRecord) ParseDates() (start, end time.Time, err error) {
	start, err = time.Parse(DateTimeFormat, or.OutlierPeriodStart)

	if err != nil {
		return
	}
	end, err = time.Parse(DateTimeFormat, or.OutlierPeriodEnd)
	return
}
