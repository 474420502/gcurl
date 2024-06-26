package gcurl

import (
	"log"
	"testing"
)

func TestNewT(t *testing.T) {
	main1()
}

func main1() {

	// cmdStr := `curl 'https://www.xxxxx.com/api-hk/heartbeat' \
	// --data '{\"name\": \"Alice\"}' \
	// -H 'accept: application/json, text/plain, */*' \
	// -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	// -H 'origin: https://www.xxxxx.com' \
	// -H 'referer: https://www.xxxxx.com/' \
	// -H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	// -H 'sec-ch-ua-mobile: ?0' \
	// -H 'sec-ch-ua-platform: "Windows"' \
	// -H 'sec-fetch-dest: empty' \
	// -H 'sec-fetch-mode: cors' \
	// -H 'sec-fetch-site: cross-site' \
	// -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`

	// opts := parseCurlCommandStr(cmdStr)
	// log.Println(opts)

}

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
