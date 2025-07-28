package gcurl

import (
	"log"
	"regexp"
	"testing"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func TestCaseWindows(t *testing.T) {

	var err error

	scurl := `curl 'https://www.xxxxx.com/api-hk/heartbeat' \
	-H 'accept: application/json, text/plain, */*' \
	-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	-H 'origin: https://www.xxxxx.com' \
	-H 'referer: https://www.xxxxx.com/' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: cross-site' \
	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`
	_, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}

	scurl2 := `curl "https://xxxx/api-hk/heartbeat" ^
	-H "accept: application/json, text/plain, */*" ^
	-H "accept-language: zh-CN,zh;q=0.9,en;q=0.8" ^
	-H "origin: https://www.xxxx.com" ^
	-H "referer: https://www.xxxx.com/" ^
	-H ^"sec-ch-ua: ^\^"Google Chrome^\^";v=^\^"123^\^", ^\^"Not:A-Brand^\^";v=^\^"8^\^", ^\^"Chromium^\^";v=^\^"123^\^"^" ^
	-H "sec-ch-ua-mobile: ?0" ^
	-H ^"sec-ch-ua-platform: ^\^"Windows^\^"^" ^
	-H "sec-fetch-dest: empty" ^
	-H "sec-fetch-mode: cors" ^
	-H "sec-fetch-site: cross-site" ^
	-H "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"`
	_, err = ParseCmd(scurl2)
	if err != nil {
		t.Error(err)
	}

	// -H ^" ^\n ^\^" ^%^"
	scurl3 := `curl "https://www.jianshu.com/shakespeare/notes/94447551/included_collections?page=1&count=7" ^
	-H "accept: application/json" ^
	-H "accept-language: zh-CN,zh;q=0.9,en;q=0.8" ^
	-H ^"cookie: locale=zh-CN; Hm_lvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712547029; _ga=GA1.2.273290375.1712547029; _ga_Y1EKTCT110=GS1.2.1712567179.2.0.1712567179.0.0.0; read_mode=day; default_font=font2; Hm_lpvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712567196; signin_redirect=https://www.jianshu.com/p/99941d7b8368; sensorsdata2015jssdkcross=^%^7B^%^22distinct_id^%^22^%^3A^%^2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c^%^22^%^2C^%^22first_id^%^22^%^3A^%^22^%^22^%^2C^%^22props^%^22^%^3A^%^7B^%^22^%^24latest_traffic_source_type^%^22^%^3A^%^22^%^E7^%^9B^%^B4^%^E6^%^8E^%^A5^%^E6^%^B5^%^81^%^E9^%^87^%^8F^%^22^%^2C^%^22^%^24latest_search_keyword^%^22^%^3A^%^22^%^E6^%^9C^%^AA^%^E5^%^8F^%^96^%^E5^%^88^%^B0^%^E5^%^80^%^BC_^%^E7^%^9B^%^B4^%^E6^%^8E^%^A5^%^E6^%^89^%^93^%^E5^%^BC^%^80^%^22^%^2C^%^22^%^24latest_referrer^%^22^%^3A^%^22^%^22^%^7D^%^2C^%^22^%^24device_id^%^22^%^3A^%^2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c^%^22^%^7D^" ^
	-H ^"if-none-match: W/^\^"c76a1ce3db1e3d9de4516a5cd05b8f6f^\^"^" ^
	-H "referer: https://www.jianshu.com/p/99941d7b8368" ^
	-H ^"sec-ch-ua: ^\^"Google Chrome^\^";v=^\^"123^\^", ^\^"Not:A-Brand^\^";v=^\^"8^\^", ^\^"Chromium^\^";v=^\^"123^\^"^" ^
	-H "sec-ch-ua-mobile: ?0" ^
	-H ^"sec-ch-ua-platform: ^\^"Windows^\^"^" ^
	-H "sec-fetch-dest: empty" ^
	-H "sec-fetch-mode: cors" ^
	-H "sec-fetch-site: same-origin" ^
	-H "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"`
	_, err = Parse(scurl3)
	if err != nil {
		t.Error(err)
	}

}

func TestMethod(t *testing.T) {
	var scurl string
	scurl = `curl -X PUT "http://httpbin.org/put"`
	_curl, err := Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	_, err = _curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X HEAD "http://httpbin.org/head"`
	_curl, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	_, err = _curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X patch "http://httpbin.org/patch"`
	_curl, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	_, err = _curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X options "http://httpbin.org/options"`
	_curl, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	_, err = _curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X DELETE "http://httpbin.org/DELETE"`
	_curl, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	_, err = _curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl "http://httpbin.org/uuid"`
	resp, err := Execute(scurl)
	if err != nil {
		t.Error(err)
	}

	if !regexp.MustCompile("uuid").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}
}

func TestParseCURL(t *testing.T) {

	scurls := []string{
		`curl 'https://saluton.cizion.com/livere' -H 'Referer: http://www.yxdm.tv/resource/9135.html' -H 'Origin: http://www.yxdm.tv' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36' -H 'Content-Type: application/json' --data-binary '{"type":"livere_pv","action":"loading","extra":{"useEagerLoading":false},"title":"哥布林杀手无删减版无暗牧无圣光 - 百度云网盘 - 全集动画下载 - 怡萱动漫","url":"http://www.yxdm.tv/resource/9135.html","consumer_seq":1020,"livere_seq":38141,"livere_referer":"www.yxdm.tv/resource/9135.html","sender":"tower","uuid":"e6213a42-41d0-4637-ad52-ccb48ba9cef1"}' --compressed`,
		`curl 'https://saluton.cizion.com/livere' -X OPTIONS -H 'Access-Control-Request-Method: POST' -H 'Origin: http://www.yxdm.tv' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36' -H 'Access-Control-Request-Headers: content-type' --compressed`,
		`curl 'https://www.google-analytics.com/r/collect' --socks5 http:127.0.0.1:7070 -H 'Referer: https://stackoverflow.com/questions/42754307/how-to-unescape-quoted-octal-strings-in-golang' -H 'Origin: https://stackoverflow.com' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36' -H 'Content-Type: text/plain;charset=UTF-8' --data-binary 'v=1&_v=j72&a=1104564653&t=pageview&_s=1&dl=https%3A%2F%2Fstackoverflow.com%2Fquestions%2F42754307%2Fhow-to-unescape-quoted-octal-strings-in-golang&ul=en-us&de=UTF-8&dt=go%20-%20How%20to%20unescape%20quoted%20octal%20strings%20in%20Golang%3F%20-%20Stack%20Overflow&sd=24-bit&sr=1476x830&vp=1412x268&je=0&_u=QACAAEAB~&jid=1066047028&gjid=2019145233&cid=572307198.1525508485&tid=UA-108242619-1&_gid=1131483813.1542548817&_r=1&cd2=%7Cstring%7Cgo%7Cescaping%7C&cd3=Questions%2FShow&z=26020125' --compressed`,
		`curl 'https://www.baidu.com/s?wd=ExpandEnv&rsv_spt=1&rsv_iqid=0xc222c428000016de&issp=1&f=8&rsv_bp=0&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_n=2&rsv_sug3=1&rsv_sug1=1&rsv_sug7=100&rsv_sug2=0&inputT=330&rsv_sug4=331' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' -H 'Upgrade-Insecure-Requests: 1' -H 'User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept-Language: zh' -H 'Cookie: BIDUPSID=88B7FC40D50C2F811E57590167144216; BAIDUID=D2066189021D32D6C36CAB19E9160526:FG=1; PSTM=1533032566; BDUSS=UNQT1ZkZW1NSzc0VmdacFowSktScWdPN2NTT3ZGTzdVMTBSaG9FMjFMSWQwNWhiQVFBQUFBJCQAAAAAAAAAAAEAAABgEGEMNDc0NDIwNTAyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB1GcVsdRnFbT; BD_UPN=123353; MCITY=-257%3A; delPer=0; BD_CK_SAM=1; PSINO=6; BD_HOME=1; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; locale=zh; H_PS_PSSID=1452_21082_18559_27401_26350_22160; sugstore=1; H_PS_645EC=78b8xypsezNZBdukGC%2Fhg6hxwjU6OnG%2BEOSA7%2BRZkLnJydHkWLS0dtpQlG6NKpJ0L8NT; BDSVRTM=0' --compressed`,
		`curl 'https://stats.g.doubleclick.net/j/collect?t=dc&aip=1&_r=3&v=1&_v=j72&tid=UA-108242619-1&cid=271874387.1533111004&jid=2011203704&gjid=1070480086&_gid=115399732.1542609235&_u=SACAAEAAEAAAAC~&z=480262738' -X POST -H 'Referer: https://stackoverflow.com/questions/28262376/parse-cookie-string-in-golang' -H 'Origin: https://stackoverflow.com' -H 'User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'Content-Type: text/plain' --compressed`,
		`curl 'https://www.google.com.hk/gen_204?s=webhp&t=aft&atyp=csi&ei=irnzW97xEIzqwQOTpYmABQ&rt=wsrt.818,aft.105,prt.105' -H 'origin: https://www.google.com.hk' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'ping-from: https://www.google.com.hk/webhp?gws_rd=cr,ssl' -H 'cookie: 1P_JAR=2018-11-20-07; NID=146=Iqtbc8EtJC9VhWWZEOMkEzscxK670vybRRaLSgEKwtJPiCaC_lRabSbBv1KWr6S3-pZ1S-VrZL4Efbwby65hjCB6SClVV7Lt0wpilw3Hr7_Uc5pzkkZOGDhVSobcl95Hs7HhuU6vb097Llu1g23NAU7mDLUB3FopfIq6lY4FpJoNhsi6L9nAnGdlXZI' -H 'x-client-data: CI22yQEIpLbJAQipncoBCKijygEY+aXKAQ==' -H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'content-type: text/ping' -H 'accept: */*' -H 'cache-control: max-age=0' -H 'authority: www.google.com.hk' -H 'ping-to: javascript:void(0);' --data-binary 'PING' --compressed`,
	}
	// Access-Control-Request-Method 方法告诉 --data-binary 默认是POST

	for _, scurl := range scurls {
		curl, err := Parse(scurl)
		if err != nil {
			t.Error(err)
		}

		if curl.Method == "" {
			t.Error("curl.Method is nil")
		}

	}
}

func TestCurlTimeout(t *testing.T) {
	scurl := `curl 'https://javtc123test.com/' --connect-timeout 1 -H 'authority: appgrowing.cn' -H 'cache-control: max-age=0' -H 'upgrade-insecure-requests: 1' -H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl, err := Parse(scurl)
	if err != nil {
		t.Error(err)
	}

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	_, err = wf.Execute()
	if err == nil {
		t.Error("not timeout")
	}
}

func TestCurlWordWrap(t *testing.T) {
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 5 
	-H 'authority: appgrowing.cn'
	-H 'cache-control: max-age=0'
	-H 'upgrade-insecure-requests: 1'
	-H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1'
	-H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8'
	-H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl, err := Parse(scurl)
	if err != nil {
		t.Error(err)
	}

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	resp, err := wf.Execute()
	if err != nil {
		t.Error(string(resp.Content()))
	}

	if len(curl.Cookies) != 3 {
		t.Error(curl.Cookies)
	}

	if len(curl.Header) != 10 { // Content-Type Cookie 不会被单独提取出来, 也是Header一种.
		t.Error(len(curl.Header), curl.Header)
	}

	scurl = `curl --header "Authorization: Bearer token_with_space" http://httpbin.org/bearer`
	curl, err = Parse(scurl)
	if err != nil {
		t.Error(err)
	}

	ses = curl.CreateSession()
	wf = curl.CreateTemporary(ses)
	resp, err = wf.Execute()
	if err != nil {
		t.Error(string(resp.Content()))
	}
	rjson := resp.Json()
	if rjson.Get("authenticated").Bool() != true {
		t.Error(`Get("authenticated").Bool() != true `)
		return
	}
	if rjson.Get("token").String() != "token_with_space" {
		t.Error(`Get("token").String() != "token_with_space"`)
		return
	}

}

func TestCurlTabCase(t *testing.T) {
	scurl := `curl 
	'http://httpbin.org/' 
	-H 'Connection: keep-alive' 
	-H 'Cache-Control: max-age=0' 
	-H 'Upgrade-Insecure-Requests: 1' 
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36' 
	-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' 
	-H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9' 
	--compressed --insecure`
	curl, err := Parse(scurl)
	if err != nil {
		t.Error(err)
	}

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	resp, err := wf.Execute()
	if err != nil {
		t.Error(string(resp.Content()))
	}

	if len(curl.Header) != 7 { // Content-Type Cookie 会被单独提取出来, 也是Header一种.
		t.Error(len(curl.Header), curl.Header)
	}

}

func TestPostFile(t *testing.T) {
	surl := `curl -X POST "http://httpbin.org/post" --data "@./tests/postfile.txt"`
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	resp, err := curl.CreateTemporary(curl.CreateSession()).Execute()
	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile("hello kids").Match(resp.Content()) {
		t.Error(resp.ContentString())
	}
}

func TestCurlPaserHttp(t *testing.T) {
	surl := ` http://httpbin.org/get -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	resp, err := curl.CreateTemporary(curl.CreateSession()).Execute()
	if err != nil {
		t.Error(err)
	}

	if !regexp.MustCompile("Accept-Encoding").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	if !regexp.MustCompile("Accept-Language").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	// log.Println(resp.Json())
	// resp.Json()
}

func TestCurlPaserHttpBody(t *testing.T) {
	surl := ` http://0.0.0.0/get/body-compressed  -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp := curl.CreateTemporary(curl.CreateSession())
	resp, err := tp.TestExecute(gserver)
	if err != nil {
		t.Error(err)
	}

	if string(resp.Content()) != "hello compress" {
		t.Error(resp.Content())
	}

	// resp.Json()
}

func TestCurlErrorCase1(t *testing.T) {
	xxxxapi := `curl  'http://httpbin.org/post' \
	-H 'authority: api.xxxx.tv' \
	-H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36' \
	-H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary3bCA1lzvhj4kBR4Q' \
	-H 'accept: */*' \
	-H 'origin: https://www.xxxx.tv' \
	-H 'sec-fetch-site: same-site' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-dest: empty' \
	-H 'referer: https://www.xxxx.tv/lives' \
	-H 'accept-language: zh-CN,zh;q=0.9' \
	--data-binary $'------WebKitFormBoundary3bCA1lzvhj4kBR4Q\r\nContent-Disposition: form-data; name="keyType"\r\n\r\n0\r\n------WebKitFormBoundary3bCA1lzvhj4kBR4Q\r\nContent-Disposition: form-data; name="body"\r\n\r\n{"deviceType":7,"requestSource":"WEB","iNetType":5}\r\n------WebKitFormBoundary3bCA1lzvhj4kBR4Q--\r\n' \
	--compressed`

	curl, err := Parse(xxxxapi)
	if err != nil {
		t.Error(err)
	}

	resp, err := curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	// t.Log(curl.String())

	if !regexp.MustCompile("keyType").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}
}

func TestCharFile(t *testing.T) {
	surl := `curl -X POST  'http://httpbin.org/post' --data-binary "@./tests/postfile.txt" `
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp := curl.CreateTemporary(nil)
	resp, err := tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile("hello kids").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	surl = `curl -X POST  'http://httpbin.org/post' --data-urlencode a=12&b=21 `
	curl, err = Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp = curl.CreateTemporary(nil)
	resp, err = tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile(`"a": "12"`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	if !regexp.MustCompile(`"b": "21"`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	surl = `curl -X POST  'http://httpbin.org/post' --data-raw a=12&b=aax `
	curl, err = Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp = curl.CreateTemporary(nil)
	resp, err = tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile(`"a": "12"`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	if !regexp.MustCompile(`"b": "aax"`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

}

func TestUser(t *testing.T) {
	surl := `curl   'http://httpbin.org/basic-auth/eson/1234567' --user eson:1234567 `
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp := curl.CreateTemporary(nil)
	resp, err := tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile(`"authenticated": true`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	surl = `curl -X POST  'http://httpbin.org/post' --user-agent golang-gcurl `
	curl, err = Parse(surl)
	if err != nil {
		t.Error(err)
	}
	tp = curl.CreateTemporary(nil)
	resp, err = tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile(`"User-Agent": "golang-gcurl"`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

}

func TestReadmeEg1(t *testing.T) {
	surl := ` http://httpbin.org/get  -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := Parse(surl)
	if err != nil {
		t.Error(err)
	}
	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)

	_, err = tp.Execute()
	if err != nil {
		log.Panic(err)
	}

}

func TestReadmeEg2(t *testing.T) {
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 1 
	-H 'authority: appgrowing.cn'
	-H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl, err := Parse(scurl)
	if err != nil {
		t.Error(err)
	}
	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	// log.Println(ses.GetCookies(wf.ParsedURL))
	// [_ga=GA1.2.1371058419.1533104518 _gid=GA1.2.896241740.1543307916 _gat_gtag_UA_4002880_19=1]
	resp, err := wf.Execute()
	if err != nil {
		log.Panic(string(resp.Content()))
	}

}

func TestReadmeEg3(t *testing.T) {
	c, err := Parse(`curl -X GET "http://httpbin.org/anything/1" -H "accept: application/json"`)
	if err != nil {
		t.Error(err)
	}
	tp := c.Temporary()
	pp := tp.PathParam(`anything/(\d+)`)
	pp.IntSet(100)
	resp, err := tp.Execute()
	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile("http://httpbin.org/anything/100").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}
}

func TestCaseLimit(t *testing.T) {

	c, err := Parse(`curl --limit-rate 200K -O http://httpbin.org/anything/100`)
	if err != nil {
		t.Error(err)
	}
	tp := c.Temporary()

	_, err = tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

}

// --abstract-unix-socket <path> Connect via abstract Unix domain socket
// --alt-svc <file name> Enable alt-svc with this cache file
// --anyauth            Pick any authentication method
// -a, --append             Append to target file when uploading
// --aws-sigv4 <provider1[:provider2[:region[:service]]]> Use AWS V4 signature authentication
// --basic              Use HTTP Basic Authentication
// --cacert <file>      CA certificate to verify peer against
// --capath <dir>       CA directory to verify peer against
// -E, --cert <certificate[:password]> Client certificate file and password
// --cert-status        Verify the status of the server cert via OCSP-staple
// --cert-type <type>   Certificate type (DER/PEM/ENG/P12)
// --ciphers <list of ciphers> SSL ciphers to use
// --compressed         Request compressed response
// --compressed-ssh     Enable SSH compression
// -K, --config <file>      Read config from a file
// --connect-timeout <fractional seconds> Maximum time allowed for connection
// --connect-to <HOST1:PORT1:HOST2:PORT2> Connect to host
// -C, --continue-at <offset> Resumed transfer offset
// -b, --cookie <data|filename> Send cookies from string/file
// -c, --cookie-jar <filename> Write cookies to <filename> after operation
// --create-dirs        Create necessary local directory hierarchy
// --create-file-mode <mode> File mode for created files
// --crlf               Convert LF to CRLF in upload
// --crlfile <file>     Use this CRL list
// --curves <algorithm list> (EC) TLS key exchange algorithm(s) to request
// -d, --data <data>        HTTP POST data
// --data-ascii <data>  HTTP POST ASCII data
// --data-binary <data> HTTP POST binary data
// --data-raw <data>    HTTP POST data, '@' allowed
// --data-urlencode <data> HTTP POST data URL encoded
// --delegation <LEVEL> GSS-API delegation permission
// --digest             Use HTTP Digest Authentication
// -q, --disable            Disable .curlrc
// --disable-eprt       Inhibit using EPRT or LPRT
// --disable-epsv       Inhibit using EPSV
// --disallow-username-in-url Disallow username in URL
// --dns-interface <interface> Interface to use for DNS requests
// --dns-ipv4-addr <address> IPv4 address to use for DNS requests
// --dns-ipv6-addr <address> IPv6 address to use for DNS requests
// --dns-servers <addresses> DNS server addrs to use
// --doh-cert-status    Verify the status of the DoH server cert via OCSP-staple
// --doh-insecure       Allow insecure DoH server connections
// --doh-url <URL>      Resolve host names over DoH
// -D, --dump-header <filename> Write the received headers to <filename>
// --egd-file <file>    EGD socket path for random data
// --engine <name>      Crypto engine to use
// --etag-compare <file> Pass an ETag from a file as a custom header
// --etag-save <file>   Parse ETag from a request and save it to a file
// --expect100-timeout <seconds> How long to wait for 100-continue
// -f, --fail               Fail fast with no output on HTTP errors
// --fail-early         Fail on first transfer error, do not continue
// --fail-with-body     Fail on HTTP errors but save the body
// --false-start        Enable TLS False Start
// -F, --form <name=content> Specify multipart MIME data
// --form-escape        Escape multipart form field/file names using backslash
// --form-string <name=string> Specify multipart MIME data
// --ftp-account <data> Account data string
// --ftp-alternative-to-user <command> String to replace USER [name]
// --ftp-create-dirs    Create the remote dirs if not present
// --ftp-method <method> Control CWD usage
// --ftp-pasv           Use PASV/EPSV instead of PORT
// -P, --ftp-port <address> Use PORT instead of PASV
// --ftp-pret           Send PRET before PASV
// --ftp-skip-pasv-ip   Skip the IP address for PASV
// --ftp-ssl-ccc        Send CCC after authenticating
// --ftp-ssl-ccc-mode <active/passive> Set CCC mode
// --ftp-ssl-control    Require SSL/TLS for FTP login, clear for transfer
// -G, --get                Put the post data in the URL and use GET
// -g, --globoff            Disable URL sequences and ranges using {} and []
// --happy-eyeballs-timeout-ms <milliseconds> Time for IPv6 before trying IPv4
// --haproxy-protocol   Send HAProxy PROXY protocol v1 header
// -I, --head               Show document info only
// -H, --header <header/@file> Pass custom header(s) to server
// -h, --help <category>    Get help for commands
// --hostpubmd5 <md5>   Acceptable MD5 hash of the host public key
// --hostpubsha256 <sha256> Acceptable SHA256 hash of the host public key
// --hsts <file name>   Enable HSTS with this cache file
// --http0.9            Allow HTTP 0.9 responses
// -0, --http1.0            Use HTTP 1.0
// --http1.1            Use HTTP 1.1
// --http2              Use HTTP 2
// --http2-prior-knowledge Use HTTP 2 without HTTP/1.1 Upgrade
// --http3              Use HTTP v3
// --ignore-content-length Ignore the size of the remote resource
// -i, --include            Include protocol response headers in the output
// -k, --insecure           Allow insecure server connections
// --interface <name>   Use network INTERFACE (or address)
// -4, --ipv4               Resolve names to IPv4 addresses
// -6, --ipv6               Resolve names to IPv6 addresses
// --json <data>        HTTP POST JSON
// -j, --junk-session-cookies Ignore session cookies read from file
// --keepalive-time <seconds> Interval time for keepalive probes
// --key <key>          Private key file name
// --key-type <type>    Private key file type (DER/PEM/ENG)
// --krb <level>        Enable Kerberos with security <level>
// --libcurl <file>     Dump libcurl equivalent code of this command line
// --limit-rate <speed> Limit transfer speed to RATE
// -l, --list-only          List only mode
// --local-port <num/range> Force use of RANGE for local port numbers
// -L, --location           Follow redirects
// --location-trusted   Like --location, and send auth to other hosts
// --login-options <options> Server login options
// --mail-auth <address> Originator address of the original email
// --mail-from <address> Mail from this address
// --mail-rcpt <address> Mail to this address
// --mail-rcpt-allowfails Allow RCPT TO command to fail for some recipients
// -M, --manual             Display the full manual
// --max-filesize <bytes> Maximum file size to download
// --max-redirs <num>   Maximum number of redirects allowed
// -m, --max-time <fractional seconds> Maximum time allowed for transfer
// --metalink           Process given URLs as metalink XML file
// --negotiate          Use HTTP Negotiate (SPNEGO) authentication
// -n, --netrc              Must read .netrc for user name and password
// --netrc-file <filename> Specify FILE for netrc
// --netrc-optional     Use either .netrc or URL
// -:, --next               Make next URL use its separate set of options
// --no-alpn            Disable the ALPN TLS extension
// -N, --no-buffer          Disable buffering of the output stream
// --no-clobber         Do not overwrite files that already exist
// --no-keepalive       Disable TCP keepalive on the connection
// --no-npn             Disable the NPN TLS extension
// --no-progress-meter  Do not show the progress meter
// --no-sessionid       Disable SSL session-ID reusing
// --noproxy <no-proxy-list> List of hosts which do not use proxy
// --ntlm               Use HTTP NTLM authentication
// --ntlm-wb            Use HTTP NTLM authentication with winbind
// --oauth2-bearer <token> OAuth 2 Bearer Token
// -o, --output <file>      Write to file instead of stdout
// --output-dir <dir>   Directory to save files in
// -Z, --parallel           Perform transfers in parallel
// --parallel-immediate Do not wait for multiplexing (with --parallel)
// --parallel-max <num> Maximum concurrency for parallel transfers
// --pass <phrase>      Pass phrase for the private key
// --path-as-is         Do not squash .. sequences in URL path
// --pinnedpubkey <hashes> FILE/HASHES Public key to verify peer against
// --post301            Do not switch to GET after following a 301
// --post302            Do not switch to GET after following a 302
// --post303            Do not switch to GET after following a 303
// --preproxy [protocol://]host[:port] Use this proxy first
// -#, --progress-bar       Display transfer progress as a bar
// --proto <protocols>  Enable/disable PROTOCOLS
// --proto-default <protocol> Use PROTOCOL for any URL missing a scheme
// --proto-redir <protocols> Enable/disable PROTOCOLS on redirect
// -x, --proxy [protocol://]host[:port] Use this proxy
// --proxy-anyauth      Pick any proxy authentication method
// --proxy-basic        Use Basic authentication on the proxy
// --proxy-cacert <file> CA certificate to verify peer against for proxy
// --proxy-capath <dir> CA directory to verify peer against for proxy
// --proxy-cert <cert[:passwd]> Set client certificate for proxy
// --proxy-cert-type <type> Client certificate type for HTTPS proxy
// --proxy-ciphers <list> SSL ciphers to use for proxy
// --proxy-crlfile <file> Set a CRL list for proxy
// --proxy-digest       Use Digest authentication on the proxy
// --proxy-header <header/@file> Pass custom header(s) to proxy
// --proxy-insecure     Do HTTPS proxy connections without verifying the proxy
// --proxy-key <key>    Private key for HTTPS proxy
// --proxy-key-type <type> Private key file type for proxy
// --proxy-negotiate    Use HTTP Negotiate (SPNEGO) authentication on the proxy
// --proxy-ntlm         Use NTLM authentication on the proxy
// --proxy-pass <phrase> Pass phrase for the private key for HTTPS proxy
// --proxy-pinnedpubkey <hashes> FILE/HASHES public key to verify proxy with
// --proxy-service-name <name> SPNEGO proxy service name
// --proxy-ssl-allow-beast Allow security flaw for interop for HTTPS proxy
// --proxy-ssl-auto-client-cert Use auto client certificate for proxy (Schannel)
// --proxy-tls13-ciphers <ciphersuite list> TLS 1.3 proxy cipher suites
// --proxy-tlsauthtype <type> TLS authentication type for HTTPS proxy
// --proxy-tlspassword <string> TLS password for HTTPS proxy
// --proxy-tlsuser <name> TLS username for HTTPS proxy
// --proxy-tlsv1        Use TLSv1 for HTTPS proxy
// -U, --proxy-user <user:password> Proxy user and password
// --proxy1.0 <host[:port]> Use HTTP/1.0 proxy on given port
// -p, --proxytunnel        Operate through an HTTP proxy tunnel (using CONNECT)
// --pubkey <key>       SSH Public key file name
// -Q, --quote <command>    Send command(s) to server before transfer
// --random-file <file> File for reading random data from
// -r, --range <range>      Retrieve only the bytes within RANGE
// --rate <max request rate> Request rate for serial transfers
// --raw                Do HTTP "raw"; no transfer decoding
// -e, --referer <URL>      Referrer URL
// -J, --remote-header-name Use the header-provided filename
// -O, --remote-name        Write output to a file named as the remote file
// --remote-name-all    Use the remote file name for all URLs
// -R, --remote-time        Set the remote file's time on the local output
// --remove-on-error    Remove output file on errors
// -X, --request <method>   Specify request method to use
// --request-target <path> Specify the target for this request
// --resolve <[+]host:port:addr[,addr]...> Resolve the host+port to this address
// --retry <num>        Retry request if transient problems occur
// --retry-all-errors   Retry all errors (use with --retry)
// --retry-connrefused  Retry on connection refused (use with --retry)
// --retry-delay <seconds> Wait time between retries
// --retry-max-time <seconds> Retry only within this period
// --sasl-authzid <identity> Identity for SASL PLAIN authentication
// --sasl-ir            Enable initial response in SASL authentication
// --service-name <name> SPNEGO service name
// -S, --show-error         Show error even when -s is used
// -s, --silent             Silent mode
// --socks4 <host[:port]> SOCKS4 proxy on given host + port
// --socks4a <host[:port]> SOCKS4a proxy on given host + port
// --socks5 <host[:port]> SOCKS5 proxy on given host + port
// --socks5-basic       Enable username/password auth for SOCKS5 proxies
// --socks5-gssapi      Enable GSS-API auth for SOCKS5 proxies
// --socks5-gssapi-nec  Compatibility with NEC SOCKS5 server
// --socks5-gssapi-service <name> SOCKS5 proxy service name for GSS-API
// --socks5-hostname <host[:port]> SOCKS5 proxy, pass host name to proxy
// -Y, --speed-limit <speed> Stop transfers slower than this
// -y, --speed-time <seconds> Trigger 'speed-limit' abort after this time
// --ssl                Try SSL/TLS
// --ssl-allow-beast    Allow security flaw to improve interop
// --ssl-auto-client-cert Use auto client certificate (Schannel)
// --ssl-no-revoke      Disable cert revocation checks (Schannel)
// --ssl-reqd           Require SSL/TLS
// --ssl-revoke-best-effort Ignore missing/offline cert CRL dist points
// -2, --sslv2              Use SSLv2
// -3, --sslv3              Use SSLv3
// --stderr <file>      Where to redirect stderr
// --styled-output      Enable styled output for HTTP headers
// --suppress-connect-headers Suppress proxy CONNECT response headers
// --tcp-fastopen       Use TCP Fast Open
// --tcp-nodelay        Use the TCP_NODELAY option
// -t, --telnet-option <opt=val> Set telnet option
// --tftp-blksize <value> Set TFTP BLKSIZE option
// --tftp-no-options    Do not send any TFTP options
// -z, --time-cond <time>   Transfer based on a time condition
// --tls-max <VERSION>  Set maximum allowed TLS version
// --tls13-ciphers <ciphersuite list> TLS 1.3 cipher suites to use
// --tlsauthtype <type> TLS authentication type
// --tlspassword <string> TLS password
// --tlsuser <name>     TLS user name
// -1, --tlsv1              Use TLSv1.0 or greater
// --tlsv1.0            Use TLSv1.0 or greater
// --tlsv1.1            Use TLSv1.1 or greater
// --tlsv1.2            Use TLSv1.2 or greater
// --tlsv1.3            Use TLSv1.3 or greater
// --tr-encoding        Request compressed transfer encoding
// --trace <file>       Write a debug trace to FILE
// --trace-ascii <file> Like --trace, but without hex output
// --trace-time         Add time stamps to trace/verbose output
// --unix-socket <path> Connect through this Unix domain socket
// -T, --upload-file <file> Transfer local FILE to destination
// --url <url>          URL to work with
// -B, --use-ascii          Use ASCII/text transfer
// -u, --user <user:password> Server user and password
// -A, --user-agent <name>  Send User-Agent <name> to server
// -v, --verbose            Make the operation more talkative
// -V, --version            Show version number and quit
// -w, --write-out <format> Use output FORMAT after completion
// --xattr              Store metadata in extended file attributes

func TestCaseIHead(t *testing.T) {

	c, err := Parse(`curl -I http://httpbin.org/anything/100`)
	if err != nil {
		t.Error(err)
	}
	tp := c.Temporary()

	resp, err := tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	// log.Println(resp.ContentString())

	// if resp.ContentString() != "" {
	// 	t.Error(`resp.ContentString() != ""`)
	// 	return
	// }
	if resp.GetResponse().Header["Content-Type"][0] != "application/json" {
		t.Error("")
		return
	}

}
