package tests

import (
	"testing"

	"github.com/474420502/gcurl"
)

func TestCase1(t *testing.T) {
	scurl4 := `	curl 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap' \
	-H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"' \
	-H 'Referer: https://github.com/474420502/gcurl/issues/3' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"'`
	c, err := gcurl.Parse(scurl4)
	if err != nil {
		t.Error(err)
	}
	_, err = c.Temporary().Execute()
	if err != nil {
		t.Error(err)
	}

}
