package prediction

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
	scriptPath = path.Join(path.Dir(filename), "fitReader.R")
}

func getBikes(stationID uint, updatetime int, weather int, temperature float64) (int, error) {
	cmd := exec.Command(
		"Rscript",
		scriptPath,
		strconv.FormatUint(uint64(stationID), 10),
		strconv.FormatBool(true), // predictBikes variable
		strconv.Itoa(updatetime),
		strconv.Itoa(weather),
		strconv.FormatFloat(temperature, 'f', 2, 64),
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
	var bikes int
	fmt.Fscan(&out, &a, &bikes)

	return bikes, nil
}
