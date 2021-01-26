package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// DetectOutliersHandler return outliers detection result or DataSet graph
func DetectOutliersHandler(w http.ResponseWriter, r *http.Request) {
	var siteID string
	var graph bool

	if siteIDParam, ok := r.URL.Query()["siteId"]; ok && len(siteIDParam[0]) > 0 {
		siteID = siteIDParam[0]
	} else {
		WriteResponse(w, 400, "Miss request param", errors.New("Expected siteId param"))
		return
	}

	if graphParam, ok := r.URL.Query()["graph"]; ok && len(graphParam) > 0 {
		graph, _ = strconv.ParseBool(graphParam[0])
	}

	ds, err := GetDataSetBySiteID(siteID)

	if err != nil {
		WriteResponse(w, 404, "Error get DataSet", err)
		return
	}
	ds.GenerateData()

	if graph {
		pl, err := MakeGraph(ds)

		if err != nil {
			WriteResponse(w, 500, "Error create graph", err)
			return
		}
		wr, err := pl.WriterTo(1024, 400, "jpeg")

		if err != nil {
			WriteResponse(w, 500, "Error create graph image", err)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		wr.WriteTo(w)
		return
	}
	results := ds.DetectOutliers()
	body, err := json.Marshal(results)

	if err != nil {
		WriteResponse(w, 500, "Error encode results", err)
		return
	}
	SetHeaders(w)
	w.Write(body)
}

func init() {
	http.HandleFunc("/api/detect_outliers", DetectOutliersHandler)
}
