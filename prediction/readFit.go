package prediction

import (
	"bytes"
	"fmt"
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

func getProbability(stationID uint, updatetime int, weather int, temperature float64) (float64, error) {
	currentTime := time.Now().Unix()

	if currentTime+shortTermThreshold >= int64(updatetime) {
		return getLongTerm(stationID, updatetime, weather)
	}

	return getShortTerm(stationID, updatetime, weather, int(currentTime))
}

func getShortTerm(stationID uint, updatetime, weather, currentTime int) (float64, error) {
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
		return 0, fmt.Errorf("%v: %s", err, errOut.String())
	}

	var a string
	var probability float64
	fmt.Fscan(&out, &a, &probability)

	return probability, nil
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
		return 0, fmt.Errorf("%v: %s", err, errOut.String())
	}

	var a string
	var probability float64
	fmt.Fscan(&out, &a, &probability)

	return probability, nil
}
