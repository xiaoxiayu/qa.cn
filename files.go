package main

import (
	"fmt"
	"strings"

	"encoding/json"

	"net/http"
	//	"time"

	"foxitsoftware.cn/quality_control/foxitqa.cn/database"
	"foxitsoftware.cn/quality_control/foxitqa.cn/file_client"
	"github.com/gin-gonic/gin"
)

type FilesHander struct {
	db          database.DBHander
	files_count int
}

func (this *FilesHander) Files(c *gin.Context) {
	this.db = database.DBHander{}
	this.db.Init("10.103.2.166")
	defer this.db.Close()

	files_count, err := this.db.CountTestFiles()
	if err != nil {
		files_count = 0
	}
	this.files_count = files_count
	//	db.CountTestFiles()
	c.HTML(http.StatusOK, "files.html", gin.H{
		"title": "Foxit QA:testfiles",
	})
}

type fileInfoData struct {
	FileName  string `json:FileName`
	StorePath string `json:StorePath`
	Size      string `json:Size`
	FileType  string `json:FileType`
	Info      string `json:Info`
	Date      string `json:Date`
}

type dbRetData struct {
	TimeStamp string         `json:timestamp`
	Data      []fileInfoData `json:data`
}

func (this *FilesHander) Upload(c *gin.Context) {
	fileclient := file_client.Server{"http://10.103.2.166:9091", 30000000, 60000000}
	//	fileclient := file_client.Server{"http://127.0.0.1:9091", 30000000, 60000000}
	filename := c.PostForm("FileName")
	info := c.PostForm("FileInfo")
	folderpath := c.PostForm("FolderPath")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{
			"ret":    2,
			"fileid": "0",
		})
		return
	}

	//	t := time.Now()
	//	store_path := c.Request.RemoteAddr + "-" + t.Format("2006010215")
	//store_path = "./Public"
	// 0: normal
	// 1: File Size is 0
	// 2: Exist in database
	// 3: Store path already exist.
	// 4: Final check file state or block error.
	// -1: internal error
	fmt.Println(folderpath + filename)
	fmt.Println(info)
	//	c.JSON(200, gin.H{
	//		"ret":  0,
	//		"info": "sd",
	//	})
	//	return
	ret, res := fileclient.Upload(file, folderpath+filename, "", info, true)
	fmt.Println("=================RET===========", ret, string(res))
	//	if ret != 0 {
	//		c.JSON(500, gin.H{
	//			"ret":  ret,
	//			"info": res,
	//		})
	//		return
	//	}
	c.JSON(200, gin.H{
		"ret":  ret,
		"info": res,
	})
	//	ret := 0
	//	res := "asdf"

}

type DataUpdateDataList struct {
	Info string
}

type DbUpdateInfoData struct {
	Data []DataUpdateDataList `json:data`
}

func (this *FilesHander) UpdateDBInfo(c *gin.Context) {

	this.db = database.DBHander{}
	this.db.Init("10.103.2.166")
	defer this.db.Close()

	FileID := c.PostForm("FileID")
	Info := c.PostForm("Info")
	where_s := fmt.Sprintf(`FileID='%s'`, FileID)
	db_data, err := this.db.Select("Info", where_s, "1")
	if err != nil {
		fmt.Println(err.Error())
	}
	//	fmt.Println("****", db_data)
	var s1 DbUpdateInfoData
	err = json.Unmarshal([]byte(db_data), &s1)
	if len(s1.Data) == 0 {
		c.JSON(500, gin.H{
			"ret":    3,
			"affect": 0,
		})
		return
	}
	all_update_info := ""
	if s1.Data[0].Info == "" {
		all_update_info = Info
	} else {
		all_update_info = s1.Data[0].Info + ";" + Info
	}

	affect := this.db.UpdateTestFileInfo(FileID, all_update_info)

	c.JSON(200, gin.H{
		"ret":    0,
		"affect": affect,
	})
}

func (this *FilesHander) createWhere(search_s string) string {
	where_i := 0
	where_s := ""
	var search_d fileInfoData
	json.Unmarshal([]byte(search_s), &search_d)
	if search_d.FileName != "" {
		//		where_s += fmt.Sprintf(`FileName='%s'`, search_d.FileName)
		//		where_i++
		mul_search := strings.Split(search_d.FileName, "|")
		if len(mul_search) > 1 {
			for _, ft := range mul_search {
				if where_i > 0 {
					where_s += fmt.Sprintf(` OR FileName='%s'`, ft)
				} else {
					where_s += fmt.Sprintf(` FileName='%s'`, ft)
				}
				where_i++
			}
		} else {
			where_s += fmt.Sprintf(`FileName='%s'`, search_d.FileName)
			where_i++
		}
	}
	if search_d.Size != "" {
		if strings.Index(search_d.Size, "=") == -1 &&
			strings.Index(search_d.Size, ">") == -1 &&
			strings.Index(search_d.Size, "<") == -1 {
			search_d.Size = "=" + search_d.Size
		}
		if where_i > 0 {
			where_s += fmt.Sprintf(` AND Size %s`, search_d.Size)
		} else {
			where_s += fmt.Sprintf(`Size %s`, search_d.Size)
		}
		where_i++
	}
	if search_d.FileType != "" {
		mul_search := strings.Split(search_d.FileType, "|")
		if len(mul_search) > 1 {
			for _, ft := range mul_search {
				if where_i > 0 {
					where_s += fmt.Sprintf(` OR FileType='%s'`, ft)
				} else {
					where_s += fmt.Sprintf(` FileType='%s'`, ft)
				}
				where_i++
			}
		} else {
			if where_i > 0 {
				where_s += fmt.Sprintf(` AND FileType='%s'`, search_d.FileType)
			} else {
				where_s += fmt.Sprintf(` FileType='%s'`, search_d.FileType)
			}
			where_i++
		}

	}
	if search_d.Date != "" {
		// >= '2015-07-22' | <= '2017-01-01'
		mul_search := strings.Split(search_d.Date, "|")
		if len(mul_search) > 1 {
			for _, ft := range mul_search {
				if where_i > 0 {
					where_s += fmt.Sprintf(` AND %s`, ft)
				} else {
					where_s += fmt.Sprintf(` Date %s`, ft)
				}
				where_i++
			}
		} else {
			if where_i > 0 {
				where_s += fmt.Sprintf(` AND %s`, search_d.Date)
			} else {
				where_s += fmt.Sprintf(` %s`, search_d.Date)
			}
			where_i++
		}

	}
	if search_d.Info != "" {
		equ_flag := "="
		if strings.Index(search_d.Info, "%") != -1 {
			equ_flag = " like "
		}
		if where_i > 0 {
			where_s += fmt.Sprintf(` AND Info%s'%s'`, equ_flag, search_d.Info)
		} else {
			where_s += fmt.Sprintf(` Info%s'%s'`, equ_flag, search_d.Info)
		}
		where_i++
	}
	if where_i > 0 {
		files_count, err := this.db.CountTestFiles_Where(where_s)
		if err != nil {
			files_count = 0
		}
		this.files_count = files_count
	}
	return where_s
}

type respondData struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  string `json:"size"`
	Date  string `json:"date"`
	Info  string `json:"info"`
	Tools string `json:"tools"`
}

func (this *FilesHander) GetData(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	this.db = database.DBHander{}
	this.db.Init("10.103.2.166")
	defer this.db.Close()

	files_count, err := this.db.CountTestFiles()
	if err != nil {
		files_count = 0
	}
	this.files_count = files_count

	search_s := c.Query("search[value]")
	where_s := ""
	fmt.Println(search_s)
	if search_s != "" {
		where_s = this.createWhere(search_s)
	}

	limit_s := c.Query("start") + "," + c.Query("length")
	//limit_s := c.Query("start") + ",10"
	//	fmt.Println(limit_s)
	//	fmt.Println(search_s)
	//	fmt.Println(where_s)
	fmt.Println(where_s)
	fmt.Println(limit_s)
	fmt.Println("search:" + search_s)

	db_data, err := this.db.Select("FileName,StorePath,FileType,Size,Info,Date", where_s, limit_s)
	if err != nil {
		c.String(http.StatusOK, "ERROR")
		return
	}

	var data dbRetData
	json.Unmarshal([]byte(db_data), &data)

	alldata := []respondData{}
	//alldata := [][]string{}
	for _, file_data := range data.Data {
		repaths := strings.SplitN(file_data.StorePath, "/mnt/mfs/Public/", -1)
		re_path := file_data.StorePath
		if len(repaths) > 1 {
			re_path = repaths[1]
		}

		//		data0 := []string{fmt.Sprintf(`<a href="http://10.103.2.166:9090/%s/%s"> %s</a>`,
		//			re_path, file_data.FileName, file_data.FileName),
		//			file_data.FileType,
		//			file_data.Size,
		//			file_data.Date,
		//			file_data.Info,
		//			`| <a href="#" data-toggle="popover" title="Preview"><span class="glyphicon glyphicon-eye-open"></span></a> |`}
		//		alldata = append(alldata, data0)

		resData := respondData{Name: fmt.Sprintf(`<a href="http://10.103.2.166:9090/%s/%s"> %s</a>`,
			re_path, file_data.FileName, file_data.FileName),
			Type:  file_data.FileType,
			Size:  file_data.Size,
			Date:  file_data.Date,
			Info:  file_data.Info,
			Tools: ` <a href="#"></a> `}
		alldata = append(alldata, resData)
	}
	//	this.files_count = 10
	c.JSON(200, gin.H{
		"draw":            c.Query("draw"),
		"recordsTotal":    this.files_count,
		"recordsFiltered": this.files_count,
		"data":            alldata,
	})
}
