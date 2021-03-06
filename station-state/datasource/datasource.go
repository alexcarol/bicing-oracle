package datasource

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BicingDataProvider provides data for the current state from the bicing API
type BicingDataProvider func() ([]byte, error)

type fakeData struct{}

// ProvideFakeData provides data that can be used when testing the app
func ProvideFakeData() ([]byte, error) {
	updateTime := time.Now().Unix()
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <bicing_stations>
     <updatetime><![CDATA[%d]]></updatetime>
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
                             </bicing_stations>`, updateTime)), nil
}

// ProvideAPIData gives you bicing data extracted from the api
func ProvideAPIData() ([]byte, error) {
	response, err := http.Get("http://wservice.viabicing.cat/v1/getstations.php?v=1")
	if err != nil {
		return nil, fmt.Errorf("Error with the request %s", err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error with the request %s", err)
	}

	return contents, nil
}
