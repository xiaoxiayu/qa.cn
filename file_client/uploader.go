package file_client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (this *Server) SetDataBase(fileSize string, fileID string,
	info string, store_path string) (bool, string) {
	data := make(url.Values)
	//		fileID := MD5CreateFromFile(absPath)
	//		bodyWriter.WriteField("fileid", fileID)
	data["FileSize"] = []string{fileSize}
	data["FileID"] = []string{fileID}
	data["Info"] = []string{info}
	data["filename"] = []string{"/" + store_path}

	target_url := this.Url + "/database"

	//fmt.Printf("%s\n", target_url)

	res, err := http.PostForm(target_url, data)
	if err != nil {
		fmt.Println(err.Error())
		return false, "Set Database ERROR."
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return false, string(result)
	}
	if strings.Contains(string(result), "ERROR") {
		return false, string(result)
	}

	return true, ""
}

func (this *Server) CheckServerExists(fileID string) (bool, string) {
	data := make(url.Values)
	data["FileID"] = []string{fileID}
	data["option"] = []string{"check"}

	target_url := this.Url + "/database"
	fmt.Println(target_url)
	res, err := http.PostForm(target_url, data)
	if err != nil {
		fmt.Println(err.Error())
		return true, ""
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return true, ""
	}
	if strings.Contains(string(result), "*NO*") {
		return false, ""
	}

	return true, string(result)
}

func (this *Server) CheckUploadfileState(file_size int64, fileID string,
	info string, store_path string, option string, is_db bool) (int, string) {

	data := make(url.Values)
	data["filename"] = []string{"/" + store_path}
	if option == "backup" {
		data["option"] = []string{"backup_folder;"}
	}

	target_url := this.Url + "/check_uploadfile"

	//fmt.Printf("%s\n", target_url)

	res, err := http.PostForm(target_url, data)
	if err != nil {
		fmt.Println(err.Error())
		return 4, "File Check ERROR."
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		//this.DeleteWithOutDB(store_path)
		return 4, "File Check ERROR."
	}
	//fmt.Printf("Recv: %s\n", result)
	serv_size, err := strconv.ParseInt(string(result), 10, 64)
	if err != nil {
		fmt.Println(err)
		//this.DeleteWithOutDB(store_path)
		return 4, "File Check ERROR."
	}
	if serv_size != file_size {
		//this.DeleteWithOutDB(store_path)
		return 4, "Size Check ERROR."
	}
	if is_db {
		if db_ret, db_str := this.SetDataBase(string(result), fileID, info, store_path); db_ret == false {
			//this.DeleteWithOutDB(store_path)
			return 4, db_str
		}
	}
	return 0, ""
}

func (this *Server) upload_file(file multipart.File, store_path string,
	option string, info string, is_db bool) (rst int, fn_info string) {
	file_path := ""
	if store_path[len(store_path)-1] == '/' {
		store_path = store_path + filepath.Base(file_path)
	}
	if store_path[len(store_path)-1] == '.' {
		store_path = filepath.Base(file_path)
	}
	backup_option := ""
	if option == "backup" {
		backup_option = "backup_folder;"
	}
	//fmt.Printf("%s\n", store_path)
	fn_info = this.Url + "/" + store_path
	//file_size := 0 //GetFileSize(file_path)
	file_size, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		fmt.Println(err.Error())
	}

	if file_size == 0 {
		rst = 1
		fn_info = "Buffer is NULL"
		return
	}
	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		fmt.Println(err.Error())
	}

	fileID := ""
	if is_db {
		fileID, _ = MD5CreateFromFile(file)
		if ret, already_path := this.CheckServerExists(fileID); ret == true {
			rst = 2 // Exists in DB.
			if already_path == "" {
				fn_info = "ERROR:DB Check Exists Failed"
			} else {
				fn_info = "FileID:" + fileID + "    \n" + already_path
			}
			return
		}
	}

	defer func(file_size int64, fileID string, info string) {
		if rst == 2 {
			return
		}
		if ck_ret, ck_str := this.CheckUploadfileState(file_size, fileID, info, store_path, option, is_db); ck_ret != 0 {
			rst = ck_ret
			fn_info = ck_str
		} else {
			fn_info = fileID
		}

	}((int64)(file_size), fileID, info)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("option", backup_option)

	bodyWriter.WriteField("block_count", "1")

	bodyWriter.WriteField("savepath", "/"+store_path)

	bodyWriter.WriteField("block_id", "1")

	bodyWriter.WriteField("block_sum", "nosum")

	fileWriter, err := bodyWriter.CreateFormFile("filename", store_path)
	if err != nil {
		rst = -1
		return
	}

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		rst = -1
		return
	}

	_, err = io.Copy(fileWriter, bytes.NewReader(bs))
	if err != nil {
		rst = -1
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	target_url := this.Url + "/normal_upload"
	resp, err := http.Post(target_url, contentType, bodyBuf)
	if err != nil {
		rst = 3
		return
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rst = -1
		return
	}

	if strings.Contains(string(resp_body), "ERROR:Existed") {
		fn_info = fn_info + " already exist."
		rst = 3 // file exist
		return
	}
	if strings.Contains(string(resp_body), "ERROR:Block store failed") {
		fn_info = "Server store error."
		rst = 3
		return
	}
	rst = 0
	return
}

func (this *Server) Upload(file multipart.File, server_path string,
	option string, info string, is_db bool) (int, string) {
	//	log_flag := "Upload"
	//	if option == "backup" {
	//		log_flag = "Backup"
	//	}
	//	if IsDir(local_path) {
	//		files := GetFilelist(local_path)
	//		for _, file_path := range files {
	//			//			fmt.Println(server_path)
	//			//			fmt.Println(len(server_path))
	//			if server_path[len(server_path)-1] == '.' {
	//				server_path = filepath.Base(local_path)
	//			}
	//			store_path := strings.SplitN(file_path, local_path, 2)[1]
	//			store_path = "/" + server_path + "/" + strings.Replace(store_path, "\\", "/", -1)
	//			//fmt.Printf("%s\n", store_path)
	//			ret, fn_info := this.upload_file(file_path, store_path, option, info, is_db)
	//			fn_info = strings.Replace(fn_info, ":9091", ":9090", 1)
	//			if ret == 0 {
	//				fmt.Printf("%s:\"%s\"\n", log_flag, fn_info)
	//			} else if ret == 3 {
	//				fmt.Printf("%s %s\n", log_flag, fn_info)
	//			} else if ret == 4 {
	//				fmt.Printf("%s %s\n", log_flag, fn_info)
	//			} else {
	//				fmt.Printf("%s ERROR:%s\n", log_flag, fn_info)
	//			}
	//		}
	//	} else {
	ret, res := this.upload_file(file, server_path, option, info, is_db)
	//	fn_info = strings.Replace(fn_info, ":9091", ":9090", 1)
	//	if ret == 0 {
	//		fmt.Printf("%s:\"%s\"\n", log_flag, fn_info)
	//	} else if ret == 3 {
	//		fmt.Printf("%s %s\n", log_flag, fn_info)
	//	} else if ret == 4 {
	//		fmt.Printf("%s %s\n", log_flag, fn_info)
	//	} else {
	//		fmt.Printf("%s ERROR:%s\n", log_flag, fn_info)
	//	}
	//	}
	return ret, res
}
