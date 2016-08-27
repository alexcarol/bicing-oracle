package fitCalculator

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

// CalculateFit calculates the fit for a station using
func CalculateFit(stationID uint, from, to int64) error {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("Error obtaining the filename")
	}
	path := path.Join(path.Dir(filename), "fitCalculator.R")

	cmd := exec.Command(
		"Rscript",
		path,
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
		return fmt.Errorf("%v: %s", err, errOut.String())
	}

	return nil
}
