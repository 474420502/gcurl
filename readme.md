# Parse curl To golang requests(https://github.com/474420502/requests)
* Easy to transform curl bash to golang code
* requests(inherit from curl bash) can add setting(config,cookie,header) and request url by you

example1:

```go
	surl := ` http://httpbin.org/get  -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl := Parse(surl)
	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)
	log.Println(ses.GetHeader())
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
```

example2:

```go
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 1 
	-H 'authority: appgrowing.cn'
	-H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh' -H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' -H 'if-none-match: W/"5bf7a0a9-ca6"' -H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	curl := Parse(scurl)
	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	log.Println(ses.GetCookies(wf.ParsedURL))
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
```