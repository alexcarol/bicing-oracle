package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
)

func main() {
	ticker := time.NewTicker(45 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				var b bytes.Buffer
				w := gzip.NewWriter(&b)

				apiData, err := datasource.APIData()
				if err != nil {
					fmt.Println(err)
					break
				}

				w.Write(apiData)
				w.Close()

				err2 := ioutil.WriteFile(fmt.Sprintf("%d.xml.gz", int32(time.Now().Unix())), b.Bytes(), 0666)
				if err2 != nil {
					fmt.Println("Filesystem write went wrong")
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}
