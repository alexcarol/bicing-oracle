package fitCalculator

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

var scriptPath string

func init() {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("Error obtaining the filename"))
	}
	scriptPath = path.Join(path.Dir(filename), "fitCalculator.R")
}

// CalculateFit calculates the fit for a station using
func CalculateFit(stationID uint, from, to int64) error {
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
