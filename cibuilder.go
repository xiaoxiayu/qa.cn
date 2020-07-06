package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	fxqacommon "xxsoftware.cn/quality_control/xxqa.cn/common"
	"github.com/gin-gonic/gin"
)

type CIBuilderHander struct {
	buildConfig jsonCI
}

var js_wizard_template string = `

`

var smart_wizard_template string = ``

type jsonCI struct {
	CI       []teamData
	User     string
	Password string
}

type teamData struct {
	Team string
	Url  string
	Data []productData
}

type productData struct {
	Product string
	Module  []moduleData
}

type CIMapData struct {
	Key       string
	SplitFlag string
	Value     []string
}

type moduleData struct {
	Name              string
	Describe          string
	Step              stepData
	DataTransform     interface{}
	DataTransformCode string
	CIProjectName     string
	CIMap             []CIMapData
}

type stepData struct {
	One   []formData
	Two   []formData
	Three []formData
}

type formData struct {
	Type    string
	Label   string
	Data    string
	Require bool
}

func createTeamLi(i int, team_s string) (team_li string) {
	if i == 0 {
		team_li = fmt.Sprintf(`<li role="presentation" class="active"><a href="#tab-%s-content" id="tab-%s-id" role="home-tab"
                                                                  data-toggle="tab" aria-expanded="true">%s</a>
                        </li>`, team_s, team_s, team_s)
	} else {
		team_li = fmt.Sprintf(`<li role="presentation"><a href="#tab-%s-content" id="tab-%s-id" role="tab"
                                                                  data-toggle="tab" aria-expanded="false">%s</a>
                        </li>`, team_s, team_s, team_s)
	}

	return
}

func createStepFormData(team_s, product_s, module_s string, data []formData) (formstr_s string) {
	//	return
	id_s := team_s + "-" + product_s + "-" + module_s
	formstr_s = `<form class="form-horizontal form-label-left">`
	for _, _data := range data {
		formstr := `<div class="form-group">`
		label_str := `<label class="control-label col-md-3 col-sm-3 col-xs-12" for="first-fxcorebranch">%s %%s </label>`
		label_str = fmt.Sprintf(label_str, _data.Label)
		if _data.Require {
			label_str = fmt.Sprintf(label_str, `<span class="required">*</span>`)
		} else {
			label_str = fmt.Sprintf(label_str, ` `)
		}
		formstr += label_str
		formstr += `<div class="col-md-6 col-sm-6 col-xs-12">`
		if _data.Type == "Input" {
			formstr += fmt.Sprintf(`<input type="text" id="%s-%s-step" name="last-name" value="%s" class="form-control col-md-7 col-xs-12">`,
				id_s, _data.Label, _data.Data)
		} else if _data.Type == "Select" {
			formstr += fmt.Sprintf(`<select id="%s-%s-step" class="form-control" required>`,
				id_s, _data.Label)
			data_l := strings.Split(_data.Data, ";")
			for _, select_data := range data_l {
				formstr += fmt.Sprintf(`<option value="">%s</option>`,
					select_data)
			}
			formstr += `</select>`
		}
		formstr += "</div></div>"

		//formstrs = append(formstrs, formstr)
		formstr_s += formstr
	}
	formstr_s += "</form>"
	return
}

func createWizardFormTemp(i int, describe_s, term_s, product_s, module_s string) (wizardform_s string) {
	class_s := `tab-pane`
	if i == 0 {
		class_s = `tab-pane active`
	}
	wizard_head := `
    <div class="%s" id="%s-%s-%s">
	<!-- Smart Wizard -->
	<p>%s</p>
    <div id="wizard-%s-%s-%s" class="form_wizard wizard_horizontal">
                                                                <ul class="wizard_steps">
                                                                    <li>
                                                                        <a href="#%s-%s-%s-step-1">
                                                                            <span class="step_no">1</span>
                          <span class="step_descr">
                                            Step 1<br/>
                                            <small>Build Setting</small>
                                        </span>
                                                                        </a>
                                                                    </li>
                                                                    <li>
                                                                        <a href="#%s-%s-%s-step-2">
                                                                            <span class="step_no">2</span>
                          <span class="step_descr">
                                            Step 2<br/>
                                            <small>Upload Setting</small>
                                        </span>
                                                                        </a>
                                                                    </li>
                                                                    <li>
                                                                        <a href="#%s-%s-%s-step-3">
                                                                            <span class="step_no">3</span>
                          <span class="step_descr">
                                            Step 3<br/>
                                            <small>Complete</small>
                                        </span>
                                                                        </a>
                                                                    </li>

                                                                </ul>`
	wizard_head = fmt.Sprintf(wizard_head,
		class_s,
		term_s, product_s, module_s,
		describe_s,
		term_s, product_s, module_s,
		term_s, product_s, module_s,
		term_s, product_s, module_s,
		term_s, product_s, module_s)
	wizardform_s += wizard_head
	wizardform_s += fmt.Sprintf(`<div id="%s-%s-%s-step-1"> %%s </div>`, term_s, product_s, module_s)
	wizardform_s += fmt.Sprintf(`<div id="%s-%s-%s-step-2"> %%s </div>`, term_s, product_s, module_s)
	wizardform_s += fmt.Sprintf(`<div id="%s-%s-%s-step-3"> %%s </div>`, term_s, product_s, module_s)
	wizardform_s += `</div>
	<!-- End SmartWizard Content -->
	</div>`
	return
}

func wizardInsertFormData(wizard_form_s string, form_str0, form_str1, form_str2 string) (wizard_form string) {
	wizard_form = fmt.Sprintf(wizard_form_s, form_str0, form_str1, form_str2)
	return
}

func createTeamModule(module_lis, wizard_form string) (teampanel_s string) {

	team_panel_head_s := ` 
        <div id="collapseOne" class="panel-collapse collapse in" role="tabpanel"
                                         aria-labelledby="headingOne">
        	<div class="panel-body">
                <div class="x_content">
                    <div class="col-xs-1">
                        <!-- required for floating -->
                        <!-- Nav tabs -->
                        <ul class="nav nav-tabs tabs-left">
							%s
						</ul>
					</div>
                	<div class="col-xs-9">
	                	<!-- Tab panes -->
	                	<div class="tab-content">
	                    %s
						</div>
					</div>
				</div>
			</div>
		</div>
					`
	teampanel_s = fmt.Sprintf(team_panel_head_s, module_lis, wizard_form)
	return

}

func createModuleLi(i int, term_s, product_s, module_s string) (module_li_s string) {
	//	return
	if i == 0 {
		module_li_s = fmt.Sprintf(`<li class="active"><a href="#%s-%s-%s" data-toggle="tab">%s</a> </li>`,
			term_s, product_s, module_s, module_s)

	} else {
		module_li_s = fmt.Sprintf(`<li><a href="#%s-%s-%s" data-toggle="tab">%s</a> </li>`,
			term_s, product_s, module_s, module_s)
	}
	return
}

func createWizardFormJS(host, team_s, ci_url, prj_s, product_s, module_s, transftom_code string, tranform interface{}, step_data stepData) (wizardform_js string) {
	id_s := team_s + "-" + product_s + "-" + module_s
	// First step data get.
	getdata_l := []string{}
	for _, _data := range step_data.One {

		if _data.Type == "Input" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$("#%s-%s-step").val()`,
				_data.Label, id_s, _data.Label))
		} else if _data.Type == "Select" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$( "#%s-%s-step option:selected").text()`,
				_data.Label, id_s, _data.Label))
		}

	}

	for _, _data := range step_data.Two {
		if _data.Type == "Input" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$("#%s-%s-step").val()`,
				_data.Label, id_s, _data.Label))
		} else if _data.Type == "Select" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$( "#%s-%s-step option:selected" ).text()`,
				_data.Label, id_s, _data.Label))
		}
	}
	for _, _data := range step_data.Three {

		if _data.Type == "Input" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$("#%s-%s-step").val()`,
				_data.Label, id_s, _data.Label))
		} else if _data.Type == "Select" {
			getdata_l = append(getdata_l, fmt.Sprintf(`"%s":$( "#%s-%s-step option:selected" ).text()`,
				_data.Label, id_s, _data.Label))
		}
	}
	getdata_s := "{"

	for _, _data := range getdata_l {
		getdata_s += _data + ","
	}
	getdata_s = getdata_s[:len(getdata_s)-1] + "}"
	js_s := `$('#wizard-%s-%s-%s').smartWizard({
            onLeaveStep:leaveAStepCallback_%s%s%s,
            onFinish:onFinishCallback_%s%s%s
        });
		function leaveAStepCallback_%s%s%s(objs, context){
			return true;
		}
		function onFinishCallback_%s%s%s(objs, context){
			var ciurl = '%s';
			var ciprj = '%s';
			var lastbuild_id = '0';
			var ws = new WebSocket("ws://%s/ci/_data");
            ws.onopen = function (evt) {
                ws.send('{"action":"start"}');
            }
            ws.onclose = function (evt) {
                //ws = null;
				//alert('ws.close');
            }
			ws.onmessage = function (evt) {
				var res = JSON.parse(evt.data);
				if (res.status == "linked") {
					ws.send('{"action":"getlastbuild","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'"}');
				} else if (res.status == "err") {
					alert(evt.data);
				} else if (res.status == "lastbuild") {					
					var ci_data = JSON.parse(res.data);
					var info = '';
					info += 'ID:\t'+ci_data.Id+'\n';
					info += 'URL:\t' + ci_data.Url+'\n';
					info += 'Building:\t' + ci_data.Building+'\n';
					info += 'EstimatedDuration:\t' + ci_data.EstimatedDuration+'\n';
					info += 'Result:\t' + ci_data.Result+'\n';
					lastbuild_id = ci_data.Id;
					
					addBuildInfo(info, 'Data Parse');
					ws.send('{"action":"createdata",'+'"term":"%s","product":"%s","module":"%s",'+'"ciurl":"'+ciurl+'","ciprj":"'+ciprj+'","formdata":'+JSON.stringify(%s)+'}');
				} else if (res.status == "cidata") {
					var data = JSON.parse(evt.data);
					var cidata = JSON.parse(data.data);
					var transform = %s;
					for (trans_key in transform) {
            			for (key_i in transform[trans_key]){
							var a_k = Object.keys(transform[trans_key][key_i])[0];
							if (cidata[trans_key] == a_k){
                    			cidata[trans_key] = transform[trans_key][key_i][a_k];
                			}
            			}
        			}
					
					var selectionStart = $('#myModal1-statuslog')[0].selectionStart;
					var selectionEnd = $('#myModal1-statuslog')[0].selectionEnd;
					
					var info = JSON.stringify(cidata);
					
					addBuildInfo(info, 'Trigger');
					ws.send('{"action":"trigger","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'",'+'"ciparameter":'+info+'}');
				} else if (res.status == "check") {
					var url_s = '/ci/lastbuild?url='+ciurl+'&project='+ciprj;
					var build_id = '';

					var info = '';
					$.get(url_s, function (data, status) {
	                        if (data.data.Id != lastbuild_id) {
	                            info += 'BuildId: \t' + data.data.Id + '\n';
								info += 'Building: \t' + data.data.Building + '\n';
								info += 'Duration: \t' + data.data.Duration + '\n';
								lastbuild_id = data.data.Id;
								addBuildInfo(info, 'Wait');
								ws.send('{"action":"waitbuild","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'"}');
	                        }else {
								info += '*Starting...*' + build_id;
								addBuildInfo(info, 'Check');
								ws.send('{"action":"waitstart","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'"}');
							}
                    	}
           			 );
				
				} else if (res.status == "waitend") {
					var info = '';
					var url_s = '/ci/lastbuild?url='+ciurl+'&project='+ciprj;
					$.get(url_s, function (data, status) {
	                        if (data.data.Building != false) {
	                            info += 'BuildId: \t' + data.data.Id + '\n';
								info += 'Buinging: \t' + data.data.Building + '\n';
								info += 'EstimatedDuration: \t' + data.data.EstimatedDuration + '\n';
								lastbuild_id = data.data.Id;
								var time_s = new Date().format('hh:mm:ss');
								addBuildInfo(info, 'Wait:'+time_s);
								ws.send('{"action":"waitbuild","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'"}');
	                        }else {
								info += '********* \n';
								info += 'Result: \t' + data.data.Result + '\n';
								addBuildInfo(info, 'END');
								$('#myModal1-titleLabel').html('%s %s End');
								ws.send('{"action":"buildend","ciurl":"'+ciurl+'","ciprj":"'+ciprj+'"}');
							}
                    	}
           			 );
				}
			}
			ws.onerror = function (evt) {
                alert("ERROR: " + evt.data);
            }
			$('#myModal1-titleLabel').html('%s %s Building...<img src="/public/img/running.gif"></img>');
			$('#myModal1').modal('show');
			$('#myModal1-statuslog').val('*Start.*\n*Get last build info.*\n');
		}
		
		`

	tranform_json, _ := json.Marshal(tranform)
	wizardform_js = fmt.Sprintf(js_s,
		team_s, product_s, module_s,
		team_s, product_s, module_s,
		team_s, product_s, module_s,
		team_s, product_s, module_s,
		team_s, product_s, module_s,
		ci_url, prj_s,
		host,
		team_s, product_s, module_s,
		getdata_s,
		tranform_json,
		product_s, module_s,
		product_s, module_s)

	return
}

type ChociesAData struct {
	Parent xml.Name `xml:"a"`
	String []string `xml:"string"`
}

type ChoicesData struct {
	Parent    xml.Name     `xml:"choices"`
	ChocicesA ChociesAData `xml:"a"`
}

type AllowSlavesData struct {
	Parent xml.Name `xml:"allowedSlaves"`
	String []string `xml:"string"`
}

type AAData struct {
	XMLName      xml.Name        `xml:"tem"`
	Name         string          `xml:"name"`
	DefaultValue string          `xml:"defaultValue"`
	Choices      ChoicesData     `xml:"choices"`
	AllowSlaves  AllowSlavesData `xml:"allowedSlaves"`
}

type PreDefineData struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type parameterDefinitionsData struct {
	Propertyies xml.Name        `xml:"parameterDefinitions"`
	DataDefine  []PreDefineData `xml:",any"`
}
type ParametersDefinitionPropertyData struct {
	Propertyies         xml.Name                 `xml:"hudson.model.ParametersDefinitionProperty"`
	ParameterDefinition parameterDefinitionsData `xml:"parameterDefinitions"`
}
type PropertiesData struct {
	Project                      xml.Name                         `xml:"project"`
	ParametersDefinitionProperty ParametersDefinitionPropertyData `xml:"hudson.model.ParametersDefinitionProperty"`
}

type CIConfigData struct {
	Project    xml.Name       `xml:"project"`
	Properties PropertiesData `xml:"properties"`
}

func GetCIJobPredefineData(config_s string) (predefine_data map[string]string) {
	v := CIConfigData{}
	err := xml.Unmarshal([]byte(config_s), &v)
	if err != nil {
		fmt.Println(err.Error())
	}

	predefine_data = make(map[string]string)
	for _, _data := range v.Properties.ParametersDefinitionProperty.ParameterDefinition.DataDefine {
		aadata := AAData{}
		xml.Unmarshal([]byte("<tem>"+_data.Content+"</tem>"), &aadata)

		if aadata.DefaultValue != "" {
			predefine_data[aadata.Name] = aadata.DefaultValue
		} else if len(aadata.Choices.ChocicesA.String) > 0 {
			predefine_data[aadata.Name] = aadata.Choices.ChocicesA.String[0]
		} else {
			if len(aadata.AllowSlaves.String) > 0 {
				predefine_data[aadata.Name] = aadata.AllowSlaves.String[0]
			}
		}

	}
	return
}

func ReGenerateData(cidata_map []CIMapData, webdatamap, predefine_data map[string]string) {
	for k, _ := range predefine_data {
		for _, m_d := range cidata_map {
			if m_d.Key == k {
				fmt.Println(m_d.Value)
				data := ""
				if len(m_d.Value) == 1 {
					predefine_data[k] = webdatamap[m_d.Value[0]]
				} else {
					for _, v_d := range m_d.Value {
						data += webdatamap[v_d] + m_d.SplitFlag
					}
					predefine_data[k] = data[:len(data)-1]
				}
			}

		}

	}
}

func (this *CIBuilderHander) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "xx Jenkins",
	})
	return

	file, e := ioutil.ReadFile("./webdata/ci.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		c.HTML(http.StatusOK, "build.html", gin.H{
			"title":            "xx QA:Build",
			"term_tab_li":      template.HTML(`<div>Configure File Error.</div>`),
			"term_tab_content": template.HTML(`<div>Configure File Error.</div>`),
			"WizardJs":         template.JS(``),
		})
		return
	}

	err := json.Unmarshal(file, &this.buildConfig)
	if err != nil {
		c.HTML(http.StatusOK, "build.html", gin.H{
			"title":            "xx QA:Build",
			"term_tab_li":      template.HTML(`<div>Configure Data Error.</div>`),
			"term_tab_content": template.HTML(fmt.Sprintf(`<div>%s</div>`, err.Error())),
			"WizardJs":         template.JS(``),
		})
		return
	}

	term_lis := ""
	wizardform_js := ""

	pagecontent := ""
	term_tab := ""
	product_content_list := []interface{}{}

	for c_i, teamdata := range this.buildConfig.CI {

		term_lis += createTeamLi(c_i, teamdata.Team)

		if c_i == 0 {
			term_tab += fmt.Sprintf(`
			<div role="tabpanel" class="tab-pane fade active in" id="tab-%s-content"
                             aria-labelledby="tab-%s-id"> %%v 
			</div>`, teamdata.Team, teamdata.Team)
		} else {
			term_tab += fmt.Sprintf(`
			<div role="tabpanel" class="tab-pane fade" id="tab-%s-content" aria-labelledby="tab-%s-id"> %%v 
			</div>`, teamdata.Team, teamdata.Team)
		}

		product_tab := `<div class="accordion" id="accordion" role="tablist" aria-multiselectable="true">`
		for _, productdata := range teamdata.Data {

			product_tab += fmt.Sprintf(`<div class="panel">
				                                    <a class="panel-heading collapsed" role="tab" id="%s-%s-heading" data-toggle="collapse"
				                                       data-parent="#accordion" href="#%s-%s-collapse" aria-expanded="false"
				                                       aria-controls="%s-%s-collapse">
				                                        <h4 class="panel-title">%s</h4>
				                                    </a>
				                                    <div id="%s-%s-collapse" class="panel-collapse collapse" role="tabpanel"
				                                         aria-labelledby="%s-%s-heading">
				                                        <div class="panel-body">%%s</div>
													</div>
											</div>`,
				teamdata.Team, productdata.Product,
				teamdata.Team, productdata.Product,
				teamdata.Team, productdata.Product,
				productdata.Product,
				teamdata.Team, productdata.Product,
				teamdata.Team, productdata.Product)

			module_li := ""

			wizard_form := ""
			term_panel := ""
			for m_i, moduledata := range productdata.Module {
				form_str0 := createStepFormData(teamdata.Team, productdata.Product, moduledata.Name, moduledata.Step.One)
				form_str1 := createStepFormData(teamdata.Team, productdata.Product, moduledata.Name, moduledata.Step.Two)
				form_str2 := createStepFormData(teamdata.Team, productdata.Product, moduledata.Name, moduledata.Step.Three)

				module_li += createModuleLi(m_i, teamdata.Team, productdata.Product, moduledata.Name)

				wizard_form_temp := createWizardFormTemp(m_i, moduledata.Describe, teamdata.Team, productdata.Product, moduledata.Name)

				wizard_form += wizardInsertFormData(wizard_form_temp, form_str0, form_str1, form_str2)

				wizardform_js += createWizardFormJS(c.Request.Host, teamdata.Team, teamdata.Url, moduledata.CIProjectName,
					productdata.Product, moduledata.Name, moduledata.DataTransformCode, moduledata.DataTransform, moduledata.Step)
			}

			term_panel = createTeamModule(module_li, wizard_form)

			product_tab = fmt.Sprintf(product_tab, term_panel)

		}
		product_tab += `</div>`

		product_content_list = append(product_content_list, product_tab)
	}

	pagecontent = fmt.Sprintf(term_tab, product_content_list...)

	c.HTML(http.StatusOK, "build.html", gin.H{
		"title":            "xx QA:Build",
		"term_tab_li":      template.HTML(term_lis),
		"term_tab_content": template.HTML(pagecontent),
		"WizardJs":         template.JS(wizardform_js),
	})
}

func (this *CIBuilderHander) getCIMap(team_s, product_s, module_s string) (cimap_data []CIMapData) {
	for _, teamdata := range this.buildConfig.CI {
		if team_s != teamdata.Team {
			continue
		}
		for _, productdata := range teamdata.Data {
			if product_s != productdata.Product {
				continue
			}
			for _, moduledata := range productdata.Module {
				if module_s != moduledata.Name {
					continue
				}
				cimap_data = moduledata.CIMap
				return
			}
		}
	}
	return
}

type CIBuildActionsCausesData struct {
	ShortDescription string
	UserId           string
	UserName         string
	UpstreamBuild    int
	UpstreamProject  string
}

type CIBuildActionsParametersData struct {
	Name  string
	Value string
}

type CIBuildActionsData struct {
	Causes     []CIBuildActionsCausesData
	Parameters []CIBuildActionsParametersData
}

type CIBuildData struct {
	Actions           []CIBuildActionsData
	Building          bool
	Duration          int
	EstimatedDuration int
	Id                string
	Url               string
	BuiltOn           string
	Result            string
}

func (this *CIBuilderHander) getCILastBuildInfo(ci_url, ci_prj string) (cibuilddata CIBuildData) {
	res, _ := fxqacommon.HTTPGet(fmt.Sprintf("http://%s:%s@%s/job/%s/lastBuild/api/json",
		this.buildConfig.User, this.buildConfig.Password, ci_url, ci_prj))
	fmt.Println(string(res))
	json.Unmarshal(res, &cibuilddata)

	return
}

func (this *CIBuilderHander) getCILastBuildInfo_Url(ci_url string) (cibuilddata CIBuildData) {
	url_tem := strings.Split(ci_url, "http://")
	res, _ := fxqacommon.HTTPGet(fmt.Sprintf("http://%s:%s@%s/lastBuild/api/json",
		this.buildConfig.User, this.buildConfig.Password, url_tem[1]))

	json.Unmarshal(res, &cibuilddata)

	return
}

func (this *CIBuilderHander) GetLastBuild(c *gin.Context) {
	url_s := c.Query("url")
	cibuilddata := this.getCILastBuildInfo_Url(url_s)
	//	js_s, err := json.Marshal(cibuilddata)
	err_i := 0
	c.JSON(200, gin.H{
		"err":  err_i,
		"data": cibuilddata,
	})
}

func (this *CIBuilderHander) getUrlFromName(ci_name string) (ci_url string) {
	for _, teamdata := range this.buildConfig.CI {
		if ci_name == teamdata.Team {
			ci_url = teamdata.Url
			return
		}

	}
	return
}

func stringInSlice(a string, params []ciTriggerData) bool {
	for _, param := range params {
		if a == param.Name {
			return true
		}
	}
	return false
}

func (this *CIBuilderHander) CI(c *gin.Context) {
	var err_s interface{}
	ret := 0
	defer func() {
		c.JSON(200, gin.H{
			"ret":  ret,
			"data": err_s,
		})
	}()
	ciurl := c.PostForm("ciurl")

	datas_name := c.PostFormArray("name")
	datas_value := c.PostFormArray("value")

	if this.buildConfig.User == "" {

		file, err := ioutil.ReadFile("./webdata/ci.json")
		if err != nil {
			err_s = err.Error()
			ret = -1
			return
		}
		err = json.Unmarshal(file, &this.buildConfig)
		if err != nil {
			err_s = err.Error()
			ret = -1
			return
		}
	}

	url_tem := strings.Split(ciurl, "http://")
	url_s := fmt.Sprintf("http://%s:%s@%s/build?delay=0sec",
		this.buildConfig.User, this.buildConfig.Password, url_tem[1])

	res, err := fxqacommon.HTTPGet(fmt.Sprintf("http://%s:%s@%s/config.xml",
		this.buildConfig.User, this.buildConfig.Password, url_tem[1]))
	if err != nil {
		err_s = "Get ci config.xml failed"
		ret = -1
		return
	}

	predefine_data := GetCIJobPredefineData(string(res))

	cidata := ciTriggerParamterData{}

	cidata.StatusCode = "303"
	cidata.RedirectTo = "."
	for i, name := range datas_name {
		if _, ok := predefine_data[name]; ok {
			citrigger_data := ciTriggerData{}
			citrigger_data.Name = name
			citrigger_data.Value = datas_value[i]
			cidata.Parameter = append(cidata.Parameter, citrigger_data)
		}

	}

	for pre_name, pre_value := range predefine_data {
		fmt.Println(pre_name)
		if !stringInSlice(pre_name, cidata.Parameter) {
			citrigger_data := ciTriggerData{}
			citrigger_data.Name = pre_name
			citrigger_data.Value = pre_value

			fmt.Println(pre_value)
			cidata.Parameter = append(cidata.Parameter, citrigger_data)
		}
	}

	js_s, err := json.Marshal(cidata)

	if err != nil {
		err_s = "Data set error"
		ret = -1
		return
	}

	ci_data := make(url.Values)

	ci_data["json"] = []string{string(js_s)}

	_, err = fxqacommon.HTTPPost(url_s, ci_data)
	if err != nil {
		err_s = "CI trigger failed"
		ret = -1
		return
	}
	cibuilddata := this.getCILastBuildInfo_Url(ciurl)
	//last_log, _ := json.Marshal(cibuilddata)
	err_s = cibuilddata
}

func (this *CIBuilderHander) CIStop(c *gin.Context) {
	var err_s interface{}
	ret := 0
	defer func() {
		c.JSON(200, gin.H{
			"ret":  ret,
			"data": err_s,
		})
	}()
	ciurl := c.PostForm("ciurl")

	if this.buildConfig.User == "" {

		file, err := ioutil.ReadFile("./webdata/ci.json")
		if err != nil {
			err_s = err.Error()
			ret = -1
			return
		}
		err = json.Unmarshal(file, &this.buildConfig)
		if err != nil {
			err_s = err.Error()
			ret = -1
			return
		}
	}

	url_tem := strings.Split(ciurl, "http://")
	url_s := fmt.Sprintf("http://%s:%s@%s/stop",
		this.buildConfig.User, this.buildConfig.Password, url_tem[1])
	fmt.Println(url_s)
	ci_data := make(url.Values)

	res, err := fxqacommon.HTTPPost(url_s, ci_data)
	if err != nil {
		//		c.JSON(200, gin.H{
		//			"ret":  -2,
		//			"data": err.Error(),
		//		})
		fmt.Println(err.Error())
		fmt.Println(string(res))
		err_s = "Stop ERROR."
		ret = -2
		return
	}

}

func (this *CIBuilderHander) SetCIBuildData(ci_url, ci_prj, term_s, product_s, module_s string, webuidata map[string]string) (predefine_data map[string]string) {
	cidata_map := this.getCIMap(term_s, product_s, module_s)

	//fmt.Println(fmt.Sprintf("http://xiaoxia_yu:xiihxvi3@%s/job/%s/config.xml", ci_url, ci_prj))
	res, _ := fxqacommon.HTTPGet(fmt.Sprintf("http://%s:%s@%s/job/%s/config.xml",
		this.buildConfig.User, this.buildConfig.Password, ci_url, ci_prj))
	//fmt.Println(string(res))
	predefine_data = GetCIJobPredefineData(string(res))
	ReGenerateData(cidata_map, webuidata, predefine_data)

	//	fmt.Println(webdatamap)
	fmt.Println(predefine_data)

	return
}

var ci_upgrader = websocket.Upgrader{}

type ciSendData struct {
	Status string `json:"status"`
	Data   string `json:"data"`
	Errstr string `json:"err"`
}

type ciRevData struct {
	Action   string            `json:"action"`
	Ciurl    string            `json:"ciurl"`
	Ciprj    string            `json:"ciprj"`
	Term     string            `json:"term"`
	Product  string            `json:"product"`
	Module   string            `json:"module"`
	FormData map[string]string `json:"formdata"`

	CIParameter map[string]string `json:ciparameter`
}

type ciTriggerData struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
type ciTriggerParamterData struct {
	Parameter  []ciTriggerData `json:"parameter"`
	StatusCode string          `json:"statusCode"`
	RedirectTo string          `json:"redirectTo"`
}

func (this *CIBuilderHander) CIState(c *gin.Context) {
	c1, err := ci_upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c1.Close()
	for {
		var _data ciSendData

		mt, message, err := c1.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		var revData ciRevData
		err = json.Unmarshal(message, &revData)
		if err != nil {
			_data.Status = "err"
			_data.Errstr = "client data error."
			fmt.Println(string(message))
			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			break
		}

		if revData.Action == "exit" {
			err = c1.WriteMessage(mt, []byte("exited"))
			if err != nil {
				fmt.Println("write:", err)
				break
			}

		} else if revData.Action == "start" {
			_data.Status = "linked"
			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}

			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}

		} else if revData.Action == "createdata" {
			predefine_data := this.SetCIBuildData(revData.Ciurl, revData.Ciprj, revData.Term, revData.Product, revData.Module, revData.FormData)

			predef_data_js, err := json.Marshal(predefine_data)
			if err != nil {
				fmt.Println(err.Error())
			}

			_data.Status = "cidata"
			_data.Data = string(predef_data_js)

			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		} else if revData.Action == "getlastbuild" {
			_data.Status = "lastbuild"
			cibuilddata := this.getCILastBuildInfo(revData.Ciurl, revData.Ciprj)
			type _buildSendData struct {
				Id                string
				Building          bool
				Url               string
				EstimatedDuration int
				BuiltOn           string
				Result            string
			}
			senddata := _buildSendData{
				Id:                cibuilddata.Id,
				Building:          cibuilddata.Building,
				Url:               cibuilddata.Url,
				EstimatedDuration: cibuilddata.EstimatedDuration,
				BuiltOn:           cibuilddata.BuiltOn,
				Result:            cibuilddata.Result}

			data, err := json.Marshal(senddata)
			if err != nil {
				fmt.Println(err.Error())
			}
			_data.Data = string(data)
			data, err = json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		} else if revData.Action == "trigger" {
			parameter := revData.CIParameter

			param := []ciTriggerData{}
			for k, v := range parameter {
				cidata_a := ciTriggerData{}
				cidata_a.Name = k
				cidata_a.Value = v
				param = append(param, cidata_a)
			}
			cidata := ciTriggerParamterData{}
			cidata.Parameter = param
			cidata.StatusCode = "303"
			cidata.RedirectTo = "."
			js_s, err := json.Marshal(cidata)
			if err != nil {
				fmt.Println(err.Error())
			}

			ci_data := make(url.Values)

			ci_data["json"] = []string{string(js_s)}
			//			fmt.Println(ci_data)

			url_s := fmt.Sprintf("http://xiaoxia_yu:xiihxvi3@%s/job/%s/build?delay=0sec", revData.Ciurl, revData.Ciprj)
			go fxqacommon.HTTPPost(url_s, ci_data)
			_data.Status = "check"
			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		} else if revData.Action == "waitstart" {
			_data.Status = "check"
			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		} else if revData.Action == "waitbuild" {
			time.Sleep(5e9)
			_data.Status = "waitend"
			data, err := json.Marshal(_data)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = c1.WriteMessage(mt, data)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		} else if revData.Action == "buildend" {

		}

	}
}
