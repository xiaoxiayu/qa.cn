package common

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/client"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/context"
)

func ErrorExcu(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"_type":"-1","_msg":"%s"}`, err.Error()))
}
func ErrorValNone(w http.ResponseWriter) {
	fmt.Fprintf(w, "%v", `{"_type":"-1","_msg":"nil"}`)
}

func ErrorNil(w http.ResponseWriter, val interface{}) {
	if val == nil {
		fmt.Fprintf(w, "%v", `{"_type":"0", "_msg":"ok"}`)
		return
	}
	w_s := ""
	switch val.(type) {
	case int:
		w_s = `"` + strconv.Itoa(val.(int)) + `"`
	case int64:
		w_s = `"` + strconv.FormatInt(val.(int64), 10) + `"`
	case string:
		w_s = `"` + val.(string) + `"`
	case bool:
		if val.(bool) {
			w_s = "true"
		} else {
			w_s = "false"
		}
	case []string:
		w_s += "["
		for _, s := range val.([]string) {
			w_s += `"` + s + `",`
		}
		w_s += w_s[:len(w_s)-2]
		w_s += "]"
	}

	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"_type":"0", "val":%s}`, w_s))
}

func ErrorParam(w http.ResponseWriter, param string) {
	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"_type":"1", "_msg":"Param '%s' Error."}`, param))
}
func HTTPPost(url_str string, data url.Values) (string, error) {
	res, err := http.PostForm(url_str, data)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func HTTPPut(url_str string, data url.Values) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", url_str, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func HTTPGet(url_str string) ([]byte, error) {
	client := &http.Client{}

	reqest, err := http.NewRequest("GET", url_str, nil)
	if err != nil {
		return nil, err
	}

	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	response, err := client.Do(reqest) //提交
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%s ErrorCode:%d", url_str, response.StatusCode)
	}

	return body, nil
}

func HTTPDelete(url_str string) error {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url_str, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func QAFileServerUpload(url, filepath, savepath string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("filename", filepath)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	if fw, err = w.CreateFormField("option"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("cover;")); err != nil {
		return
	}
	if fw, err = w.CreateFormField("block_count"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("1")); err != nil {
		return
	}

	if fw, err = w.CreateFormField("savepath"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("/" + savepath)); err != nil {
		return
	}

	if fw, err = w.CreateFormField("block_id"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("1")); err != nil {
		return
	}

	if fw, err = w.CreateFormField("block_sum"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("nosum")); err != nil {
		return
	}

	w.Close()

	// Now that you have a form, you can submit it to your handler.
	upload_page := "/normal_upload"
	req, err := http.NewRequest("POST", url+upload_page, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Println(res.StatusCode)
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func QAFileServerGet(url, filepath string) ([]byte, error) {
	return HTTPGet(url + filepath)
}

//func ParseHashValue(hash_val string) (map[string]string, int, error) {
//	val_map := make(map[string]string)
//	vals := strings.Split(hash_val, " ")
//	val_size := len(vals)
//	if val_size%2 != 0 {
//		return nil, -1, fmt.Errorf("Values len error.")
//	}
//	val_key := ""
//	for id, val := range vals {
//		i := id % 2
//		if i == 0 {
//			val_key = val
//		} else if i == 1 {
//			val_map[val_key] = val
//		}
//	}
//	return val_map, val_size, nil
//}

func ParseHashValue(fields, vals []string) (map[string]string, error) {
	val_map := make(map[string]string)
	for i, v := range fields {
		val_map[v] = vals[i]
		i++
	}
	return val_map, nil
}

func WatcheFile(filepath string, doFunc func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					doFunc()
				}
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filepath)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func EtcdInit(etcd_endpoints string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: []string{etcd_endpoints},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	kapi := client.NewKeysAPI(c)
	return kapi, nil
}

func EtcdSet(kApi client.KeysAPI, key, value string) (*client.Response, error) {
	return kApi.Set(context.Background(), key, value, &client.SetOptions{Dir: false})
}

func ReadCfg(cfg_paths []string) interface{} {
	cfg := flag.String("cfg", "", "Configure file.")
	flag.Parse()
	if *cfg != "" {
		cfg_paths = append(cfg_paths, *cfg)
	}
	var config interface{}
	//	fmt.Println(cfg_path)
	for _, cfg_path := range cfg_paths {
		if _, err := toml.DecodeFile(cfg_path, &config); err != nil {
			//			log.Println(err)
			continue
		}

		break
	}

	return config
}

func PathExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func GetLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				//			case *net.IPAddr:
				//				ip = v.IP
			}

			if ip.String() == "127.0.0.1" ||
				ip == nil ||
				strings.Count(ip.String(), ".") != 3 {
				continue
			}

			return ip, nil
		}
	}
	return nil, fmt.Errorf("Not found IP")
}
