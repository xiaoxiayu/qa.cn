package main

import (
	"encoding/json"
	"fmt"
	"os"
	//"html/template"
	"net/http"

	"io/ioutil"
	"os/exec"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"encoding/base64"

	fxqacommon "xxsoftware.cn/quality_control/xxqa.cn/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cmap "github.com/streamrail/concurrent-map"
)

type TestState struct {
	MachineInfo cmap.ConcurrentMap
}

type stateCacheData struct {
	Val        []string `json:val`
	Ci         string   `json:ci`
	Starttime  string   `json:starttime`
	Index      string   `json:index`
	Output     string   `json:output`
	Progress   string   `json:progress`
	Creator    string   `json:creator`
	Info       string   `json:info`
	Status     string   `json:status`
	Type       string   `json:"_type"`
	Excutor    string   `json:"executor"`
	ExcutorPid string   `json:"executor_pid"`
}

func (this *TestState) State(c *gin.Context) {
	c.HTML(http.StatusOK, "state.html", gin.H{
		"title":           "xx QA:state of the test run",
		"status_trs_html": "NODATA",
		"fxcore_host":     "ws://" + c.Request.Host + "/test/state/_data",
	})
}

var upgrader = websocket.Upgrader{}

type retData struct {
	Type       string `json:"type"`
	Buf        string `json:"buf"`
	TaskName   string `json:"task"`
	Progress   string `json:"progress"`
	CreateUrl  string `json:"create_url"`
	InfoUrl    string `json:"info_url"`
	OutputUrl  string `json:"output_url"`
	Creator    string `json:"creator"`
	Count      int    `json:"count"`
	Status     string `json:"status"`
	StatusInfo string `json:"status_info"`
	StartTime  string `json:"start_time"`
	Excutor    string `json:"executor"`
	ExcutorPid string `json:"executor_pid"`
}

func getTaskSet() (stateCacheData, error) {
	//G_FXQACACHE_IP := "http://10.103.129.80:32460/v1/platform/cache"
	G_FXQACACHE_IP := "http://10.103.129.80:32457"
	cache_task_url := G_FXQACACHE_IP + "/set?key=FXQA_TESTING"
	var task_val stateCacheData
	res, err := fxqacommon.HTTPGet(cache_task_url)
	if err != nil {
		fmt.Println(err.Error())
		return task_val, err
	}
	//fmt.Println(string(res))

	err = json.Unmarshal(res, &task_val)
	if err != nil {
		fmt.Println(err.Error())
		return task_val, nil
	}
	return task_val, nil
}

func (this *TestState) StateData(c *gin.Context) {
	//G_FXQACACHE_IP := "http://10.103.129.80:32460/v1/platform/cache"
	G_FXQACACHE_IP := "http://10.103.129.80:32457"

	//	origin := c.Request.Header.Get("Origin")
	//	fmt.Println(origin)
	//	whiteList := "domain"
	c.Request.Header.Del("Origin")
	//	if origin == whiteList {
	//		c.Request.Header.Del("Origin")
	//		fmt.Println("===========")
	//	}
	//	fmt.Println(c.Request.Header)
	c1, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	//	fmt.Println("99999")
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c1.Close()
	for {
		mt, message, err := c1.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}

		if string(message) == "createalldata" ||
			string(message) == "updatealldata" ||
			string(message) == "doupdate" {
			task_val, err := getTaskSet()
			if err != nil {
				fmt.Println("eeeeeee")
				break
			}
			if string(message) == "updatealldata" {
				err = c1.WriteMessage(mt, []byte(fmt.Sprintf(`{"type":"updatealldata_start","count":%d}`, len(task_val.Val))))
				if err != nil {
					fmt.Println("write:", err)
					break
				}
			} else {
				for _, task_name := range task_val.Val {

					cache_data_url := G_FXQACACHE_IP + "/hash?key=FXQA-" + task_name

					res, err := fxqacommon.HTTPGet(cache_data_url)
					if err != nil {
						fmt.Println(err.Error())
						fmt.Println("aafff")
					}

					var task_detail stateCacheData
					err = json.Unmarshal(res, &task_detail)
					if err != nil {
						fmt.Println(err.Error())
						fmt.Println("oooopkpk")
					}
					if task_detail.Info == "" {
						task_detail.Info = "Running"
					}
					if task_detail.Creator == "" {
						task_detail.Creator = "#"
					}

					if task_detail.Status == "" {
						task_detail.Status = "Unknown"
					}
					if task_detail.Progress == "" {
						task_detail.Progress = "0"
					}

					if task_detail.Type == "3" {
						task_detail.Status = "Stopped"
					}

					var _data retData
					_data.TaskName = task_name
					_data.CreateUrl = task_detail.Ci
					_data.InfoUrl = task_detail.Index
					_data.Creator = task_detail.Creator
					_data.OutputUrl = task_detail.Output
					_data.Progress = task_detail.Progress
					_data.Status = task_detail.Status
					_data.StatusInfo = task_detail.Info
					_data.StartTime = task_detail.Starttime
					_data.Excutor = task_detail.Excutor
					_data.ExcutorPid = task_detail.ExcutorPid
					_data.Type = "create"
					if string(message) == "doupdate" {
						_data.Type = "update"
					}

					data, err := json.Marshal(_data)
					if err != nil {
						fmt.Println(err.Error())
						fmt.Println("wefwef")
					}

					err = c1.WriteMessage(mt, data)
					if err != nil {
						fmt.Println("write:", err)
						break
					}
				}
			}
			if string(message) == "createalldata" {
				err = c1.WriteMessage(mt, []byte(fmt.Sprintf(`{"type":"createall_end","count":%d}`, len(task_val.Val))))
				if err != nil {
					fmt.Println("write:", err)
					break
				}
			}
			if string(message) == "doupdate" {
				err = c1.WriteMessage(mt, []byte(fmt.Sprintf(`{"type":"updateall_end","count":%d}`, len(task_val.Val))))
				if err != nil {
					fmt.Println("write:", err)
					break
				}
			}
			if string(message) == "end" {
				return
			}
			//}
		}
	}
}

type HeartbeatData struct {
	MachineName string `json:"machine"`
	MachineInfo string `json:"machineinfo"`
	Status      string `json:"status"` // free or busy
	TestInfo    string `json:"testinfo"`
	TimeStamp   int64  `json:"time"`
}

func (this *TestState) TimerCleaner() {
	var timeout_define time.Duration = 15
	var ticker *time.Ticker = time.NewTicker(timeout_define * time.Second)
	for t := range ticker.C {
		//fmt.Println("Tick at", t.Unix())
		for info := range this.MachineInfo.Items() {
			//fmt.Println("Info:", info)
			temdata, _ := this.MachineInfo.Get(info)
			data := temdata.(*HeartbeatData)
			if (t.Unix() - data.TimeStamp) > int64(timeout_define) {
				this.MachineInfo.Remove(info)
				fmt.Println("Remove:", info)
			}
		}
	}
}

func (this *TestState) HeartBeat(c *gin.Context) {
	cdata := c.PostForm("data")
	var data HeartbeatData
	json.Unmarshal([]byte(cdata), &data)
	fmt.Println(data)
	data.TimeStamp = time.Now().Unix()
	this.MachineInfo.Set(data.MachineName, &data)
}

type AllMachineInfo struct {
	Info []HeartbeatData
}

type TestInfoData struct {
	Name       string `json:"name"`
	Info       string `json:"info"`
	Creator    string `json:"creator"`
	CreateTime string `json:"create_time"`
}

type EchoTestInfoData struct {
	MachineName string `json:"machine"`
	MachineInfo string `json:"machineinfo"`
	TestName    string `json:"testname"`
	TestInfo    string `json:"testinfo"`
	UpdateTime  string `json:"time"`
	Status      string `json:"status"`
	CreateTime  string `json:"create_time"`
	Creator     string `json:"creator"`
}

func (this *TestState) GetHeartBeatInfo(c *gin.Context) {
	var allInfo []EchoTestInfoData //AllMachineInfo
	for info := range this.MachineInfo.Items() {
		//fmt.Println("Info:", info)
		temdata, _ := this.MachineInfo.Get(info)
		data := temdata.(*HeartbeatData)

		var testinfo TestInfoData
		json.Unmarshal([]byte(data.TestInfo), &testinfo)

		tm := time.Unix(data.TimeStamp, 0)

		echoData := EchoTestInfoData{
			MachineName: data.MachineName,
			MachineInfo: data.MachineInfo,
			Status:      data.Status,
			TestName:    testinfo.Name,
			TestInfo:    testinfo.Info,
			UpdateTime:  tm.String(),
			CreateTime:  testinfo.CreateTime,
			Creator:     testinfo.Creator,
		}
		allInfo = append(allInfo, echoData)
	}
	c.JSON(200, gin.H{
		"ret":  0,
		"data": allInfo,
	})
}

type EchoMachineStatusData struct {
	MachineName string `json:"machine"`
	Status      string `json:"status"`
}

func (this *TestState) GetStatus(c *gin.Context) {
	status_s := c.Query("status")

	var allInfo []EchoMachineStatusData
	for info := range this.MachineInfo.Items() {
		//fmt.Println("Info:", info)
		temdata, _ := this.MachineInfo.Get(info)
		data := temdata.(*HeartbeatData)
		if status_s != "" {
			if data.Status != status_s {
				continue
			}
		}

		var testinfo TestInfoData
		json.Unmarshal([]byte(data.TestInfo), &testinfo)

		echoData := EchoMachineStatusData{
			MachineName: data.MachineName,
			Status:      data.Status,
		}
		allInfo = append(allInfo, echoData)
	}
	c.JSON(200, gin.H{
		"ret":  0,
		"data": allInfo,
	})
}

func (this *TestState) SetGUISummary(c *gin.Context) {

	action := c.PostForm("action")
	if action == "" {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": "action empty",
		})
		return
	}

	tablename := "gui" //c.PostForm("testgroup")

	db, err := sql.Open("sqlite3", "./webdata/guitest_summary.db")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()
	if action == "addgroup" {

		sqlStmt := fmt.Sprintf(`
	create table %s (testname text not null primary key, updatetime text, testindex text, ci text, status text, result text, resulturl text);
	`, tablename)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
	} else if action == "update" {
		testname := c.PostForm("testname")
		testversion := c.PostForm("version")
		ci := c.PostForm("ci")
		testindex := c.PostForm("testindex")
		teststatus := c.PostForm("status")
		testresult := c.PostForm("result")
		result_url := c.PostForm("resulturl")

		table_line_name := testname + "_" + testversion
		if testversion == "" {
			table_line_name = testname
		}

		stmt, err := db.Prepare(fmt.Sprintf("select testname from %s where testname = ?", tablename))
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow(table_line_name).Scan(&name)
		if err != nil {

		}
		if name != "" {
			ret, err := db.Exec(fmt.Sprintf("update %s set updatetime='%s', testindex='%s', ci='%s', status='%s', result='%s', resulturl='%s' where testname='%s'",
				tablename,
				time.Now().Format("2006-01-02 15:04:05"),
				testindex,
				ci,
				teststatus,
				testresult,
				result_url,
				table_line_name))
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			affectlined, _ := ret.RowsAffected()
			if affectlined == 0 {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": "insert failed",
				})
				return
			}

			//tx.Commit()
		} else {

			tx, err := db.Begin()
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			stmt, err = tx.Prepare(fmt.Sprintf(
				"insert into %s(testname, updatetime, testindex, ci, status, result, resulturl) values(?, ?, ?, ?, ?, ?, ?)",
				tablename))
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})

				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(table_line_name,
				time.Now().Format("2006-01-02 15:04:05"),
				testindex,
				ci,
				teststatus,
				testresult,
				result_url)
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			tx.Commit()
		}
	}
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

type testSummaryData struct {
	Testname   string
	Updatetime string
	Testindex  string
	Ci         string
	Status     string
	Result     string
	Resulturl  string
}

func (this *TestState) GetGUISummary(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./webdata/guitest_summary.db")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()

	rows, err := db.Query("select testname, updatetime, testindex, ci, status, result, resulturl from gui")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer rows.Close()
	data := []testSummaryData{}
	for rows.Next() {
		_testData := testSummaryData{}
		err = rows.Scan(&_testData.Testname,
			&_testData.Updatetime,
			&_testData.Testindex,
			&_testData.Ci,
			&_testData.Status,
			&_testData.Result,
			&_testData.Resulturl)
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
		data = append(data, _testData)

	}
	err = rows.Err()
	if err != nil {
		c.JSON(200, gin.H{
			"ret":  -2,
			"msg":  err.Error(),
			"data": []string{},
		})
		return
	}

	c.JSON(200, gin.H{
		"ret":  0,
		"data": data,
	})
}

func (this *TestState) DelGUISummary(c *gin.Context) {
	testname := c.Param("testname_version")

	db, err := sql.Open("sqlite3", "./webdata/guitest_summary.db")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("delete from gui where testname='%s'", testname))
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func (this *TestState) SetGUIMachineMsg(c *gin.Context) {
	//	x, _ := ioutil.ReadAll(c.Request.Body)
	action := c.PostForm("action")
	if action == "" {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": "action empty",
		})
		return
	}

	tablename := "gui" //c.PostForm("testgroup")

	db, err := sql.Open("sqlite3", "./webdata/guitest_machine_message.db")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()
	if action == "addgroup" {

		sqlStmt := fmt.Sprintf(`
	create table %s (ip text not null primary key, message text, updatetime text, purpose text);
	`, tablename)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
	} else if action == "update" {
		ip_s := c.PostForm("ip")
		msg_s := c.PostForm("msg")
		purpose := c.PostForm("purpose")

		stmt, err := db.Prepare(fmt.Sprintf("select ip from %s where ip = ?", tablename))
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow(ip_s).Scan(&name)
		if err != nil {

		}
		if name != "" {
			ret, err := db.Exec(fmt.Sprintf("update %s set updatetime='%s', message='%s',purpose='%s' where ip='%s'",
				tablename,
				time.Now().Format("2006-01-02 15:04:05"),
				msg_s,
				purpose,
				ip_s))
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			affectlined, _ := ret.RowsAffected()
			if affectlined == 0 {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": "insert failed",
				})
				return
			}

			//tx.Commit()
		} else {

			tx, err := db.Begin()
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			stmt, err = tx.Prepare(fmt.Sprintf(
				"insert into %s(ip, message, updatetime, purpose) values(?, ?, ?, ?)",
				tablename))
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})

				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(ip_s,
				msg_s,
				time.Now().Format("2006-01-02 15:04:05"),
				purpose)
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}
			tx.Commit()
		}
	}
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

type machineMsgData struct {
	Ip         string
	UpdateTime string
	Purpose    string
	Message    string
}

func (this *TestState) GetGUIMachineMsg(c *gin.Context) {
	ip_s := c.Query("ip")

	tablename := "gui" //c.PostForm("testgroup")
	db, err := sql.Open("sqlite3", "./webdata/guitest_machine_message.db")
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()
	if ip_s == "" {
		rows, err := db.Query("select ip, message, updatetime, purpose from gui")
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
		defer rows.Close()

		machineDatas := []machineMsgData{}
		for rows.Next() {
			var Ip, Message, UpdateTime, Purpose string
			err = rows.Scan(&Ip, &Message, &UpdateTime, &Purpose)
			if err != nil {
				c.JSON(200, gin.H{
					"ret": -2,
					"msg": err.Error(),
				})
				return
			}

			machineDatas = append(machineDatas,
				machineMsgData{Ip: Ip,
					Message:    Message,
					UpdateTime: UpdateTime,
					Purpose:    Purpose})
		}
		err = rows.Err()
		if err != nil {
			c.JSON(200, gin.H{
				"ret": -2,
				"msg": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"ret": 0,

			"msg": machineDatas,
		})
		return
	}

	stmt, err := db.Prepare(fmt.Sprintf("select message, updatetime, purpose from %s where ip = ?", tablename))
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	defer stmt.Close()
	var Msg, UpdateTime, Purpose string
	err = stmt.QueryRow(ip_s).Scan(&Msg, &UpdateTime, &Purpose)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"ret":        0,
		"message":    Msg,
		"updatetime": UpdateTime,
		"purpose":    Purpose,
	})
}

func (this *TestState) FuzzChartCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": "name empty",
		})
		return
	}
	plot := c.PostForm("plot")
	if plot == "" {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": "plot empty",
		})
		return
	}
	bitmap := c.PostForm("bitmap")
	if bitmap == "" {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": "bitmap empty",
		})
		return
	}
	//	data, _ := fxqacommon.HTTPGet("http://127.0.0.1:9091/fuzz-data")
	decodePlotBytes, err := base64.StdEncoding.DecodeString(plot)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err.Error())
	}
	//fmt.Println(bitmap)

	decodeBitmapBytes, err := base64.StdEncoding.DecodeString(bitmap)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err.Error())
	}
	//	json.Marshal()
	folder_path := "./frontend/dist/" + name
	os.MkdirAll(folder_path, 0700)

	folder_outputpath := "./frontend/dist/output/fuzz_state_" + name
	os.MkdirAll(folder_outputpath, 0700)

	err = ioutil.WriteFile(folder_path+"/plot_data", decodePlotBytes, 0644)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}
	err = ioutil.WriteFile(folder_path+"/fuzz_bitmap", decodeBitmapBytes, 0644)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -2,
			"msg": err.Error(),
		})
		return
	}

	//	run_cmd := []string{}
	cmd := exec.Command("./afl-plot", folder_path, folder_outputpath)
	stdout, _ := cmd.StdoutPipe()
	if err != nil {
		return
	}
	cmd.Start()

	ioutil.ReadAll(stdout)
	c.JSON(200, gin.H{
		"ret":  0,
		"data": "output/fuzz_state_" + name,
	})
}
