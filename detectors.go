package main

import (
	"errors"
	"time"
)

// ThreeSigmasOutlierDetector outlier detection by 3-sigmas rule
func ThreeSigmasOutlierDetector(ds DataSet, mv MetricValues) (*OutlierDetectOutput, error) {
	step, err := ParseDuration(ds.TimeStep)

	if err != nil {
		return nil, err
	}

	vals := mv.Values.FilterByMinValue(ds.GetMinValue())
	total := len(vals)

	if total == 0 {
		return nil, errors.New("No available values")
	}

	var x, y, n, i, j int
	start := vals[0].Date
	end := vals[total-1].Date

	out := &OutlierDetectOutput{
		SiteID:                  ds.SiteID,
		TimeAgo:                 ds.TimeAgo,
		TimeStep:                ds.TimeStep,
		OutliersDetectionMethod: ThreeSigmas,
		DateStart:               start.Format(DateTimeFormat),
		DateEnd:                 end.Format(DateTimeFormat),
		Result: OutliersDetectResult{
			Warnings: make([]OutlierDetectResultRecord, 0),
			Alarms:   make([]OutlierDetectResultRecord, 0),
		},
	}

	defer func(starTime time.Time) {
		out.CheckTimeStart = starTime.UTC().Format(DateTimeFormat)
		out.CheckTimeEnd = time.Now().UTC().Format(DateTimeFormat)
	}(time.Now())

	for ; start.Before(end); start = start.Add(step) {
		x = y

		for y < total && vals[y].Date.Day() == start.Day() {
			y++
		}

		mean, std := vals[x:y].GetMeanStDev()

		warnUpperLimit := mean + std*ds.OutliersDetection.OutliersMultipler
		alarmUpperLimit := mean + std*ds.OutliersDetection.StrongOutliersMultipler

		for i = x; i < y; i++ {

			if vals[i].Value > warnUpperLimit && x > 0 && y < total {
				n = i
				j = i

				for n >= x && vals[n].Value > mean {
					n--
				}
				for j <= y && vals[j].Value > mean {
					j++
				}

				r := OutlierDetectResultRecord{
					Metric:             mv.Metric,
					Attribute:          mv.Attribute,
					OutlierPeriodStart: vals[n].Date.Format(DateTimeFormat),
					OutlierPeriodEnd:   vals[j].Date.Format(DateTimeFormat),
				}

				if vals[i].Value > alarmUpperLimit {
					out.Result.Alarms = append(out.Result.Alarms, r)
				} else {
					out.Result.Warnings = append(out.Result.Warnings, r)
				}
				i = j
			}
		}
	}
	return out, nil
}
