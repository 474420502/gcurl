package main

import (
	"encoding/gob"
	"flag"
	"os"
)

func main() {

	// os.Stdin.WriteString("curl 'http://localhost:7070/api-hk/heartbeat' -H 'accept: application/json, text/plain, */*' -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8'")
	flag.Parse()
	args := flag.Args()
	err := gob.NewEncoder(os.Stdout).Encode(&args)
	if err != nil {
		panic(err)
	}

}
