package gcurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// 测试 -w/--write-out 功能
func TestWriteOut(t *testing.T) {
	// 启动一个简单的服务器
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	}))
	defer srv.Close()

	// 构造 curl 命令，使用 write-out 格式
	format := "HTTP:%{http_code} TIME:%{time_total}"
	sc := fmt.Sprintf("curl -w \"%s\" %s", format, srv.URL)

	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	resp, err := Execute(sc)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	out, _ := ioutil.ReadAll(r)
	output := string(out)
	if !strings.Contains(output, "HTTP:201") {
		t.Errorf("write-out output missing status code: %s", output)
	}
	if !strings.Contains(output, "TIME:") {
		t.Errorf("write-out output missing time_total: %s", output)
	}
}

// 测试 -f/--fail 功能
func TestFailOnError(t *testing.T) {
	// 返回404
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer srv.Close()

	sc := fmt.Sprintf("curl -f %s", srv.URL)
	_, err := Execute(sc)
	if err == nil {
		t.Errorf("expected error for 404 with -f, got nil")
	}
}

// 测试 -J/--remote-header-name 功能
func TestRemoteHeaderName(t *testing.T) {
	body := "filecontent"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=demo.txt")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	// 创建临时目录
	tmpdir, _ := ioutil.TempDir("", "gcurl_test")
	defer os.RemoveAll(tmpdir)

	// 解析命令
	sc := fmt.Sprintf("curl -J -o %s %s", tmpdir, srv.URL)
	curl, err := Parse(sc)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	// 执行请求
	ses := curl.CreateSession()
	req := curl.CreateTemporary(ses)
	resp, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	// 保存到文件
	err = curl.SaveToFile(resp)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// 检查文件存在且内容正确
	path := tmpdir + string(os.PathSeparator) + "demo.txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != body {
		t.Errorf("file content mismatch: got %s", string(data))
	}
}
