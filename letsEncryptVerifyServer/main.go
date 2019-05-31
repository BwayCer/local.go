package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	insCmdFlag := handleFlag()
	keyList := insCmdFlag.KeyList.Cache

	http.HandleFunc("/", func(resp http.ResponseWriter, requ *http.Request) {
		urlPath := requ.URL.Path
		key := keyList[urlPath]
		if key == "" {
			record(requ, false)
			resp.WriteHeader(http.StatusNotFound)
		} else {
			record(requ, true)
			resp.Header().Set("Content-Type", "text/plain")
			resp.Header().Set("Content-Length", strconv.Itoa(len(key)))
			resp.WriteHeader(http.StatusOK)
			io.WriteString(resp, key)
		}
	})

	panic(http.ListenAndServe("0.0.0.0:80", nil))
}

type keyList struct {
	Cache map[string]string
}

func (self keyList) String() string {
	return fmt.Sprintf("%v", self.Cache)
}

func (self keyList) Set(value string) error {
	cutList := strings.Split(value, "-=")

	if len(cutList) < 2 {
		err := fmt.Errorf("路徑與金鑰的設定值 %q 不如預期。", value)
		return err
	} else if len(cutList) != 2 {
		err := fmt.Errorf("路徑與金鑰的設定值 %q 解析到過多參數。", value)
		return err
	}

	path := cutList[0]
	key := cutList[1]

	if path == "" || key == "" {
		panic("未提供指定路徑或金鑰。")
	} else if string(path[0]) != "/" {
		panic("指定路徑必須以 \"/\" 為開頭。")
	}

	self.Cache[path] = key
	return nil
}

type CmdFlag struct {
	KeyList keyList
}

func handleFlag() CmdFlag {
	insCmdFlag := CmdFlag{KeyList: keyList{Cache: make(map[string]string)}}
	flag.Var(insCmdFlag.KeyList, "s", "設定路徑與金鑰。 (ex: `path-=key`)")
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
