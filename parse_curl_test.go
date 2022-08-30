package gcurl

import (
	"log"
	"regexp"
	"testing"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func TestMethod(t *testing.T) {
	var scurl string
	scurl = `curl -X PUT "http://httpbin.org/put"`
	_, err := Parse(scurl).CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X HEAD "http://httpbin.org/head"`
	_, err = Parse(scurl).CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X patch "http://httpbin.org/patch"`
	_, err = Parse(scurl).CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X options "http://httpbin.org/options"`
	_, err = Parse(scurl).CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl -X DELETE "http://httpbin.org/DELETE"`
	_, err = Parse(scurl).CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	scurl = `curl "http://httpbin.org/uuid"`
	resp, err := Execute(scurl)
	if err != nil {
		t.Error(err)
	}

	if !regexp.MustCompile("uuid").Match(resp.Content()) {
		t.Error(resp.Content())
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
		`curl 'https://saluton.cizion.com/livere' -F 'formkey="formvalue"'`,
	}
	// Access-Control-Request-Method 方法告诉 --data-binary 默认是POST

	for _, scurl := range scurls {
		curl := Parse(scurl)

		if curl.Method == "" {
			t.Error("curl.Method is nil")
		}

	}
}

func TestCurlTimeout(t *testing.T) {
	scurl := `curl 'https://javtc123test.com/' --connect-timeout 1 -H 'authority: appgrowing.cn' -H 'cache-control: max-age=0' -H 'upgrade-insecure-requests: 1' -H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl := Parse(scurl)

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	_, err := wf.Execute()
	if err == nil {
		t.Error("not timeout")
	}
}

func TestCurlWordWrap(t *testing.T) {
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 1 
	-H 'authority: appgrowing.cn'
	-H 'cache-control: max-age=0'
	-H 'upgrade-insecure-requests: 1'
	-H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1'
	-H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8'
	-H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl := Parse(scurl)

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	resp, err := wf.Execute()
	if err != nil {
		t.Error(string(resp.Content()))
	}

	if len(curl.Cookies) != 3 {
		t.Error(curl.Cookies)
	}

	if len(curl.Header) != 9 { // Content-Type Cookie 会被单独提取出来, 也是Header一种.
		t.Error(len(curl.Header), curl.Header)
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
	curl := Parse(scurl)

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
	curl := Parse(surl)
	resp, err := curl.CreateTemporary(curl.CreateSession()).Execute()
	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile("hello kids").Match(resp.Content()) {
		t.Error(resp.Content())
	}
}

func TestCurlPaserHttp(t *testing.T) {
	surl := ` http://httpbin.org/get  -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl := Parse(surl)
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

	curl := Parse(xxxxapi)
	resp, err := curl.CreateTemporary(nil).Execute()
	if err != nil {
		t.Error(err)
	}

	t.Log(curl.String())

	if !regexp.MustCompile("keyType").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}
}

func TestCharFile(t *testing.T) {
	surl := `curl -X POST  'http://httpbin.org/post' --data-binary "@./tests/postfile.txt" `
	curl := Parse(surl)
	tp := curl.CreateTemporary(nil)
	resp, err := tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile("hello kids").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	surl = `curl -X POST  'http://httpbin.org/post' --data-urlencode a=12&b=21 `
	curl = Parse(surl)
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
	curl = Parse(surl)
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
	curl := Parse(surl)
	tp := curl.CreateTemporary(nil)
	resp, err := tp.Execute()
	if err != nil {
		panic(err)
	}
	if !regexp.MustCompile(`"authenticated": true`).Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	surl = `curl -X POST  'http://httpbin.org/post' --user-agent golang-gcurl `
	curl = Parse(surl)
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
	curl := Parse(surl)
	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)
	// log.Println(ses.GetHeader())
	// map[Accept-Encoding:[gzip, deflate] Accept-Language:[zh-CN,zh;q=0.9] Connection:[keep-alive]]
	resp, err := tp.Execute()
	if err != nil {
		log.Panic(err)
	}

	log.Println(string(resp.Content()))
	//     ------response-----
	//     "args": {},
	//     "headers": {
	//     "Accept-Encoding": "gzip, deflate",
	//     "Accept-Language": "zh-CN,zh;q=0.9",
	//     "Connection": "keep-alive,close",
	//     "Host": "httpbin.org",
	//     "User-Agent": "Go-http-client/1.1"
	//     },
	//     "origin": "172.17.0.1",
	//     "url": "http://httpbin.org/get"
	// }
}

func TestReadmeEg2(t *testing.T) {
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 1 
	-H 'authority: appgrowing.cn'
	-H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl := Parse(scurl)
	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	// log.Println(ses.GetCookies(wf.ParsedURL))
	// [_ga=GA1.2.1371058419.1533104518 _gid=GA1.2.896241740.1543307916 _gat_gtag_UA_4002880_19=1]
	resp, err := wf.Execute()
	if err != nil {
		log.Panic(string(resp.Content()))
	}

	log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"headers": {
	// 	  "Accept-Encoding": "gzip, deflate, br",
	// 	  "Accept-Language": "zh",
	// 	  "Authority": "appgrowing.cn",
	// 	  "Connection": "close",
	// 	  "Cookie": "_ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1",
	// 	  "Host": "httpbin.org",
	// 	  "If-Modified-Since": "Fri, 23 Nov 2018 06:39:37 GMT",
	// 	  "If-None-Match": "W/\"5bf7a0a9-ca6\"",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/get"
	//   }
}

func TestReadmeEg3(t *testing.T) {
	c := Parse(`curl -X GET "http://httpbin.org/anything/1" -H "accept: application/json"`)
	tp := c.Temporary()
	pp := tp.PathParam(`anything/(\d+)`)
	pp.IntSet(100)
	resp, err := tp.Execute()
	if err != nil {
		t.Error(err)
	}
	log.Println(string(resp.Content()))
	if !regexp.MustCompile("http://httpbin.org/anything/100").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}
}
