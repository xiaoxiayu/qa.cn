// main.go
package file_client

import (
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"
	//	"fmt"
	"os"
)

type Server struct {
	Url       string
	BlockSize int64
	LargeSize int64
}

func MergeMD5(md5vec []string) []string {
	var md5Str string
	var allmd5 string
	for _, f := range md5vec {
		allmd5 += f
	}
	h := md5.New()
	h.Write([]byte(allmd5))

	md5Str = hex.EncodeToString(h.Sum(nil))

	md5vec = nil
	md5vec = append(md5vec, md5Str)
	return md5vec
}

func MD5CreateFromFile(f multipart.File) (md5s string, err error) {
	defer func() {
		_, err = f.Seek(0, os.SEEK_SET)
	}()
	buff := make([]byte, 20480)
	var md5vec []string
	for n, err := f.Read(buff); err == nil; n, err = f.Read(buff) {
		h := md5.New()
		h.Write([]byte(string(buff[:n])))

		md5vec = append(md5vec, hex.EncodeToString(h.Sum(nil)))
		if len(md5vec) >= 128 {
			md5vec = MergeMD5(md5vec)
		}
	}
	if len(md5vec) > 0 {
		md5vec = MergeMD5(md5vec)
	}
	md5s = md5vec[0]
	return
}
