package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	insCmdFlag := handleFlag()

	if insCmdFlag.Path == "" || string(insCmdFlag.Path[0]) != "/" || insCmdFlag.Key == "" {
		panic("未提供指定路徑或指定金鑰。")
	}

	server := &http.Server{Addr: "0.0.0.0:8080"}

	http.HandleFunc(insCmdFlag.Path, func(resp http.ResponseWriter, requ *http.Request) {
		record(requ, true)

		key := insCmdFlag.Key

		resp.Header().Set("Content-Type", "text/plain")
		resp.Header().Set("Content-Length", strconv.Itoa(len(key)))
		resp.WriteHeader(http.StatusOK)
		io.WriteString(resp, key)
	})

	http.HandleFunc("/", func(resp http.ResponseWriter, requ *http.Request) {
		record(requ, false)
		resp.WriteHeader(http.StatusNotFound)
	})

	panic(server.ListenAndServe())
}

type CmdFlag struct {
	Path string
	Key  string
}

func handleFlag() CmdFlag {
	insCmdFlag := CmdFlag{}
	flag.StringVar(&insCmdFlag.Path, "path", "", "指定路徑。")
	flag.StringVar(&insCmdFlag.Key, "key", "", "指定金鑰。")
	flag.Parse()
	return insCmdFlag
}

func record(requ *http.Request, ynTargetPath bool) {
	var targetPathPrompt string
	now := time.Now()
	ip, _, _ := net.SplitHostPort(requ.RemoteAddr)
	urlPath := requ.URL.Path

	if ynTargetPath {
		targetPathPrompt = " (Arrival at the destination)"
	}

	fmt.Printf("%s come form %s to %s\n", now.Format("15:04:05.000000"), ip, urlPath+targetPathPrompt)
}
