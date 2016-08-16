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
	err := readFit(1)
	fmt.Println(err)
}

func readFit(stationID uint) error {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("Error obtaining the filename")
	}
	path := path.Join(path.Dir(filename), "fitReader.R")

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

	if err != nil {
		return fmt.Errorf("%v: %s", err, errOut.String())
	}
	fmt.Println(out.String())
	//fmt.Println(errOut.String())

	return nil
}
