package main

import (
	"errors"
	"time"
)

// ThreeSigmasOutlierDetector outlier detection by 3-sigmas rule
func ThreeSigmasOutlierDetector(ds DataSet, mv MetricValues) (*OutlierDetectOutput, error) {
	checkStartTime := time.Now().UTC()

	if len(mv.Values) == 0 {
		return nil, errors.New("Empty values")
	}
	timeDuration, timeStep, err := ds.GetTimeAgoAndTimeStepDurations()

	if err != nil {
		return nil, err
	}
	endDate := mv.Values[mv.Values.Len()-1].Date.Truncate(timeStep)
	startDate := endDate.Add(-timeDuration).Round(timeStep)

	minDetectionValue := float64(timeStep / timeDuration * time.Duration(ds.MinVisitorsPerTimeStep))
	values := mv.Values.FilterByDatesAndMinValue(startDate, endDate, minDetectionValue)

	if values.Len() == 0 {
		return nil, errors.New("No available values for detecting")
	}
	output := ds.MakeOutlierOutput(startDate, endDate)

	defer func() {
		output.CheckTimeStart = checkStartTime.UTC().Format(DateTimeFormat)
		output.CheckTimeEnd = time.Now().UTC().Format(DateTimeFormat)
	}()

	parts := values.BreakIntoPieces(timeStep)
	means := make([]float64, len(parts))
	stDevs := make([]float64, len(parts))

	for i, part := range parts {
		means[i], stDevs[i] = part.GetMeanStDev()
	}
	commonMean, _ := MeanStDev(means...)
	commonStDev, _ := MeanStDev(stDevs...)

	warnUpperLimit := commonMean + commonStDev*ds.OutliersDetection.OutliersMultipler
	alarmUpperLimit := commonMean + commonStDev*ds.OutliersDetection.StrongOutliersMultipler

	for indx, part := range parts {
		if means[indx] < commonMean {
			continue
		}
		for start := 1; start < part.Len(); start++ {
			if part[start].Value <= warnUpperLimit {
				continue
			}
			stop := start

			for part[stop].Value > means[indx]+stDevs[indx] && stop < part.Len() {
				stop++
			}
			result := OutlierDetectResultRecord{
				OutlierPeriodStart: part[start-1].Date.Format(DateTimeFormat),
				OutlierPeriodEnd:   part[stop].Date.Format(DateTimeFormat),
				Metric:             mv.Metric,
				Attribute:          mv.Attribute,
			}

			if mean, _ := part[start:stop].GetMeanStDev(); mean > alarmUpperLimit {
				output.Result.Alarms = append(output.Result.Alarms, result)
			} else {
				output.Result.Warnings = append(output.Result.Warnings, result)
			}
			start = stop
		}
	}
	return output, nil
}
