package main

import (
	"fmt"
	"github.com/alexcarol/bicing-api/station-state/repository"
	"github.com/alexcarol/bicing-api/station-state/parser"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func apiFakeDataProvider() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
    <bicing_stations>
     <updatetime><![CDATA[1415996588]]></updatetime>
      <station>
        <id>1</id>
          <type>BIKE</type>
            <lat>41.397952</lat>
              <long>2.180042</long>
                <street><![CDATA[Gran Via Corts Catalanes]]></street>
                  <height>21</height>
                    <streetNumber>760</streetNumber>
                      <nearbyStationList>24, 369, 387, 426</nearbyStationList>
                        <status>OPN</status>
                          <slots>0</slots>
                            <bikes>24</bikes>
                             </station>
                             </bicing_stations>`)
}

func doAPIRequest() []byte {
	response, err := http.Get("http://wservice.viabicing.cat/v1/getstations.php?v=1")
	if err != nil {
		fmt.Printf("Error with the request %s", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error with the request %s", err)
		os.Exit(1)
	}

	return contents
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	go func(s *repository.StationStateStorage) {
		for {
			select {
			case <-ticker.C:
				data := parser.Parse(apiFakeDataProvider())
				s.PersistCollection(data)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(repository.NewStorage())

	<-quit
}
