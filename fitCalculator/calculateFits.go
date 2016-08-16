package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

func main() {
	err := calculateFit(1)
	fmt.Println(err)
}

func calculateFit(stationID uint) error {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("Error obtaining the filename")
	}
	path := path.Join(path.Dir(filename), "fitCalculator.R")

	cmd := exec.Command(
		"Rscript",
		path,
		strconv.FormatUint(uint64(stationID), 10),
	)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	fmt.Println(out.String())
	if err != nil {
		return fmt.Errorf("%v: %s", err, errOut.String())
	}

	return nil
}
