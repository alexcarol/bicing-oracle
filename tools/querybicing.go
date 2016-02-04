package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
)

func main() {
	db, dberr := sql.Open("mysql", "alex:alexpassword@/bicing_raw")
	if dberr != nil {
		panic(dberr)
	}

	result, queryErr := db.Query("show databaes")
	if queryErr != nil {
		panic(queryErr)
	}

	fmt.Println(result)

	fmt.Println("Starting querybicing")
	ticker := time.NewTicker(45 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Querying bicing")
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
