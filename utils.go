package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// WriteResponse write http response
func WriteResponse(w http.ResponseWriter, code int, message string, err error) {
	SetHeaders(w)
	w.WriteHeader(code)
	response := JSONResponse{Code: code, Message: message}

	if err == nil {
		response.Status = true
	} else {
		response.Error = err.Error()
	}

	jsResponse, _ := json.Marshal(response)
	w.Write(jsResponse)
}

// SetHeaders set default reponse headers
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
}

// ReadFile open and read local file
func ReadFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, fmt.Errorf("Error open config file: %s", err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("Error read config file: %s", err)
	}
	return body, nil
}

// StartServer start listen http server on given port
func StartServer(port uint) {
	if port == 0 || port > math.MaxUint16 {
		log.Fatalln("Wrong server port, port must be in the range from 1 to 65535")
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    HTTPReadTimeout,
		WriteTimeout:   HTTPWriteTimeout,
		MaxHeaderBytes: MaxHearedBytes,
	}
	defer server.Shutdown(context.Background())
	log.Printf("Start http server on 0.0.0.0:%d\n", port)
	log.Fatalln(server.ListenAndServe())
}

// MeanStDev calc mean and standart deviation
func MeanStDev(args ...float64) (mean float64, stdDev float64) {
	l := float64(len(args))

	if l <= 1 {
		return
	}

	var sum float64

	for i := range args {
		sum += args[i]
	}
	mean = sum / l

	for i := range args {
		stdDev += math.Pow(args[i]-mean, 2)
	}
	stdDev = math.Sqrt(stdDev / (l - 1))
	return mean, stdDev
}

// GenerateValue generates random values ​​depending on the time
// and generates outliers on the 11th at 12 and 18 hours
func GenerateValue(dt time.Time) float64 {
	hour := dt.Hour()
	var min, max float64 = 750, 975

	switch hour {
	case 0, 1, 2, 3, 4, 5, 6, 22, 23:
		min, max = 750, 900
	case 7, 8, 9, 10, 11, 12, 19, 20, 21:
		min, max = 900, 950
	case 13, 14, 15, 16, 17, 18:
		min, max = 950, 975
	}

	// Generate outliers
	if dt.Day() == 11 {
		switch hour {
		case 12:
			min, max = 2000, 2000
		case 18:
			min, max = 2450, 2450
		}
	}
	return (rand.Float64()*(max-min) + min)
}
