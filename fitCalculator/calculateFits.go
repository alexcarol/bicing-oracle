package fitCalculator

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

type fitRequest struct {
	stationID       uint
	from, to        int64
	responseChannel chan<- error
}

var calculate = make(chan fitRequest)

var defaultErrorChannel = make(chan error)

func init() {
	go func() {
		for {
			err := <-defaultErrorChannel
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("Error obtaining the filename"))
	}
	scriptPath := path.Join(path.Dir(filename), "fitCalculator.R")

	go func(scriptPath string) {
		for {
			data := <-calculate

			err := doCalculateFit(data.stationID, data.from, data.to, scriptPath)
			if err != nil {
				data.responseChannel <- fmt.Errorf(
					"Error calculating fit for %d, %d, %d in script %s : %v\n",
					data.stationID,
					data.from,
					data.to,
					scriptPath,
					err,
				)
			} else {
				data.responseChannel <- nil
			}
		}

	}(scriptPath)
}

const (
	// TODO maybe these constants should be used only internally
	DefaultFrom = 1465653992 // TODO fix
	DefaultTo   = 9999999999 // TODO fix
)

// ScheduleCalculate schedules a fit calculation for the chosen station
func ScheduleCalculate(stationID uint) {
	go func() {
		calculate <- fitRequest{stationID, DefaultFrom, DefaultTo, defaultErrorChannel}
	}()
}

// CalculateFit calculates the fit for a station
func CalculateFit(stationID uint, from, to int64) error {
	var responseChannel = make(chan error)

	calculate <- fitRequest{stationID, from, to, responseChannel}

	return <-responseChannel
}

func doCalculateFit(stationID uint, from, to int64, scriptPath string) error {
	cmd := exec.Command(
		"Rscript",
		scriptPath,
		strconv.FormatUint(uint64(stationID), 10),
		strconv.FormatInt(from, 10),
		strconv.FormatInt(to, 10),
	)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf(
			"error calculating fit: %v: out:%s, err:%s",
			err,
			out.String(),
			errOut.String(),
		)
	}

	return nil
}
