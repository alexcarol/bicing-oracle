package parser

import "testing"

func TestXMLParsesReturnsErrWhenNotAssignedParsableXML(t *testing.T) {
	_, err := ParseXML([]byte("<asdf>"))
	if err == nil {
		t.FailNow()
	}
}
