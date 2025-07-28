package gcurl

import (
	"log"
	"testing"
)

func TestIssue(t *testing.T) {

	scurl := ` 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap' \
  -H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"' \
  -H 'Referer: https://github.com/474420502/gcurl/issues/3' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"'`
	cu, err := Parse(scurl)
	if err != nil {
		log.Panic(err)
	}
	resp, err := cu.Temporary().Execute()
	if err != nil {
		log.Panic(err)
	}

	log.Println(resp.ContentString())

}
