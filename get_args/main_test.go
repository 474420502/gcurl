package main_test

import (
	"bytes"
	"encoding/gob"
	"log"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("bash", "-c", "go run main.go curl 'http://localhost:7070/api-hk/heartbeat' -H 'accept: application/json, text/plain, */*' -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8'")
	// err := cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }
	data, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	var buf = bytes.NewBuffer(data)
	var args []string
	err = gob.NewDecoder(buf).Decode(&args)
	if err != nil {
		panic(err)
	}
	log.Println(args)
}
