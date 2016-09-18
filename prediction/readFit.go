package prediction

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"
)

var longTermScriptPath string
var shortTermScriptPath string

func init() {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("Error obtaining the filename"))
	}
	longTermScriptPath = path.Join(path.Dir(filename), "fitReader.R")
	shortTermScriptPath = path.Join(path.Dir(filename), "shortTermPredictor.R")
}

const shortTermThreshold = 1800 // 30 minutes

func getProbability(stationID uint, updatetime int, weather int, temperature float64, currentBikes int) (float64, error) {
	currentTime := time.Now().Unix()

	if currentTime+shortTermThreshold <= int64(updatetime) {
		log.Println("Long term", currentTime, shortTermThreshold, updatetime)
		return getLongTerm(stationID, updatetime, weather)
	}

	log.Println("Short term", currentTime, shortTermThreshold, updatetime)

	return getShortTerm(stationID, updatetime, weather, int(currentTime), currentBikes)
}

func getShortTerm(stationID uint, updatetime, weather, currentTime int, currentBikes int) (float64, error) {
	cmd := exec.Command(
		"Rscript",
		shortTermScriptPath,
		strconv.FormatUint(uint64(stationID), 10),
		strconv.FormatBool(true), // predictBikes variable
		strconv.Itoa(updatetime),
		strconv.Itoa(currentTime),
		strconv.Itoa(weather),
		strconv.Itoa(weather),
	)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		return 0, fmt.Errorf("%v: %s, %s", err, errOut.String(), out.String())
	}

	var a int
	var futurePrediction float64
	var currentPrediction float64
	fmt.Fscan(&out, &a, &a, &futurePrediction, &currentPrediction)

	if float64(currentBikes)+futurePrediction-currentPrediction > 0 {
		return 1, nil
	}

	return 0, nil
}

func getLongTerm(stationID uint, updatetime int, weather int) (float64, error) {
	cmd := exec.Command(
		"Rscript",
		longTermScriptPath,
		strconv.FormatUint(uint64(stationID), 10),
		strconv.FormatBool(true), // predictBikes variable
		strconv.Itoa(updatetime),
		strconv.Itoa(weather),
	)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		return 0, fmt.Errorf("%v: Error: %s, Regular output: %s", err, errOut.String(), out.String())
	}

	var a string
	var probability float64
	fmt.Fscan(&out, &a, &probability)

	return probability, nil
}
