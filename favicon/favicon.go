package favicon

import (
	"io/ioutil"
	"local/mimeList"
	"log"
	"net/http"
	"path"
	"strconv"
)

type Ico struct {
	Path string
	buf  []byte
}

func (self *Ico) ReadImage() []byte {
	buf, err := ioutil.ReadFile(self.Path)
	if err != nil {
		log.Fatal(err)
	}

	self.buf = buf
	return buf
}

func (self *Ico) Server(resp http.ResponseWriter, requ *http.Request) {
	resp.Header().Set("Content-Type", mimeList.ByExt(path.Ext(self.Path)))
	resp.Header().Set("Content-Length", strconv.Itoa(len(self.buf)))
	resp.WriteHeader(http.StatusOK)
	resp.Write(self.buf)
}
