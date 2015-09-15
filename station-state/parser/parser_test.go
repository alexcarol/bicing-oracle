package parser

import "testing"

func TestXMLParsesReturnsErrWhenNotAssignedParsableXML(t *testing.T) {
	_, err := ParseXML([]byte("<asdf>"))
	if err == nil {
		t.FailNow()
	}
}

func TestParseXMLReturnsAProperCollection(t *testing.T) {
	var correctData = []byte(`<?xml version="1.0" encoding="UTF-8"?>
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

	_, err := ParseXML(correctData)
	if err == nil {
		t.FailNow()
	}

}

func assertEquals(t *testing.T) {

}
