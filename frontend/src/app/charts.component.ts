import { Component, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { MatSnackBar } from '@angular/material';
import { FileUploader, FileItem, ParsedResponseHeaders } from 'ng2-file-upload';



@Component({
    selector: 'app-charts',
    templateUrl: './charts.component.html',
    styleUrls: ['./charts.component.css']
})


export class ChartsComponent implements OnInit {
    text: string = "";
    win_ip: string = '10.103.129.202:9091';
    options: any = { maxLines: 1000, printMargin: false };
    init_option: string = "// Simple Demo.\n\
{\n\
    title: {\n\
        text: 'ChartDemo'\n\
    },\n\
    legend: {\n\
        data: ['legend0', 'legend1']\n\
    },\n\
    toolbox: {\n\
        show: true,\n\
        feature: {\n\
            dataView: { readOnly: false },\n\
            restore: {},\n\
            saveAsImage: {}\n\
        }\n\
    },\n\
    xAxis: [\n\
        {\n\
            name: 'x',\n\
            type: 'category',\n\
            data: ['data0', 'data1', 'data2']\n\
        }\n\
    ],\n\
    yAxis: [\n\
        {\n\
            name: 'y',\n\
            type: 'value'\n\
        }\n\
    ],\n\
    series: [\n\
        {\n\
            name: 'legend0',\n\
            type: 'bar',\n\
            data: [4, 5, 6]\n\
        },\n\
        {\n\
            name: 'legend1',\n\
            type: 'bar',\n\
            data: [1, 2, 3]\n\
        }\n\
    ]\n\
}";

    file_service_url: string = 'files/upload';
    public uploader: FileUploader = new FileUploader({
        url: this.file_service_url,
        method: "POST",
        itemAlias: "file",
        autoUpload: false,
        isHTML5: true
    });

    chartOption: any = {};
    logChartOptions: Array<any> = [];

    constructor(public snackBar: MatSnackBar) {
        this.uploader.onAfterAddingFile = (file) => {
            file.withCredentials = false;
            
            var reader = new FileReader();
            reader.onloadend = lfile => {
                var contents: any = lfile.target;
                // console.info(contents);
                // this.text = contents.result;
                // console.info('afteradd', lfile);
                this.CreateChart(file._file.name, contents.result); 
            };
            reader.readAsText(file._file);
        };
    }

    ngOnInit(): void {
        this.text = this.init_option;
        // eval("this.chartOption=" + this.init_option);
    }

    ngAfterViewInit() {
        eval("this.chartOption=" + this.init_option);
    }

    onChange(code) {
        try {
            eval("this.chartOption=" + code);
        } catch (e) {
            setTimeout(() => this.snackBar.open((e), 'Option code error.', { duration: 3000 }))
        }

    }

    selectedFileOnChanged(e: any): void {
        // this.file_selected = true;
        // this.readSingleFile('');
        console.info('change', e);

    }

    getInitOption(): any {
		var colors = ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff', '#222200'];

		var option = {
			//color: colors,

			title:{
				text:''
			},
			tooltip: {
				trigger: 'axis',
				axisPointer: {
					type: 'cross'
				}
			},
			grid: {
				right: '20%'
			},
			dataZoom: [
				{
					show: true,
					start: 0,
					end: 100
				},
				{
					type: 'inside',
					start: 0,
					end: 100
				},
				{
					show: true,
					yAxisIndex: [0, 1, 2],
					filterMode: 'empty',
					width: 30,
					height: '80%',
					showDataShadow: false,
					right: '1%'
				}
			],
			toolbox: {
				//orient: 'vertical',
				feature: {
					mark : {show: true},
					dataView : {show: true, readOnly: false},
					magicType: {show: true, type: ['line', 'bar', 'stack', 'tiled']},
					restore: {show: true},
					saveAsImage: {
						show: true,
						title : '保存为图片',
						type : 'jpeg',
						lang : ['点击本地保存']
					}

				}
			},
			legend: {
				data:[]
			},
			xAxis: [
				{
					type: 'category',
					axisTick: {
						alignWithLabel: true
					},
					data: []
				}
			],
			yAxis: [
				{
					type: 'value',
					name: 'time',
					min: 0,
					max: 250,
					position: 'left',
					axisLine: {
						lineStyle: {
							color: colors[0]
						}
					},
					axisLabel: {
						formatter: '{value} ms'
					}
				},
				{
					type: 'value',
					name: 'cpu',
					min: 0,
					max: 400,
					position: 'right',
					axisLine: {
						lineStyle: {
							color: colors[1]
						}
					},
					axisLabel: {
						formatter: '{value} %'
					}
				},

				{
					type: 'value',
					name: 'memory',
					min: 0,
					max: 3000,
					position: 'right',
					offset: 80,
					axisLine: {
						lineStyle: {
							color: colors[2]
						}
					},
					axisLabel: {
						formatter: '{value} mb'
					}
				},
			],
			series: [

			]
		};
		return option;
	}

    CreateChart(logname: string, contents: string) {
		// var main_chart = echarts.init(document.getElementById(element_id));

		let option = this.getInitOption();
		option['animation']	= false;

		// main_chart.showLoading();
		var lines = contents.split('\n');
		var config_data = JSON.parse(lines[0]);
		var x_asix_name = config_data.x;

		var step_len = config_data.step.length;

		// Time configure.
		var time_config = config_data.time;
		if (typeof(time_config) == 'undefined') {
			time_config = [];
		}
		var time_max = 0;
		var config_cnt = 0;
		for (var i = 0; i < time_config.length; i++) {
			option.legend.data.push('time-'+time_config[i]);
			var ser = {
				name:'time-'+time_config[i],
				stack: 'time',
				yAxisIndex: 0,
				type:'bar',
				markPoint:{
					data:[]
				},
				data:[]
			};
			option.series.push(ser);
			config_cnt++;
			if (config_cnt > 5) {
				option.legend.data.push('');
				config_cnt = 0;
			}
		}
		
		// CPU configure.
		var cpu_config = config_data.cpu;
		if (typeof(cpu_config) == 'undefined') {
			cpu_config = [];
		}
		for (var i = 0; i < cpu_config.length; i++) {
			option.legend.data.push('cpu-'+cpu_config[i]);
			var ser = {
				name:'cpu-'+cpu_config[i],
				stack: 'cpu',
				yAxisIndex: 1,
				type:'bar',
				markPoint:{
					data:[]
				},
				data:[]
			};
			option.series.push(ser);
			config_cnt++;
			if (config_cnt > 5) {
				option.legend.data.push('');
				config_cnt = 0;
			}
		}

		// Memory configure.
		var mem_max = 0;
		var mem_config = config_data.memory;
		if (typeof(mem_config) == 'undefined') {
			mem_config = [];
		}
		for (var i = 0; i < mem_config.length; i++) {
			option.legend.data.push('memory-'+mem_config[i]);
			var ser = {
				name:'memory-'+mem_config[i],
				stack: 'memory',
				yAxisIndex: 2,
				type:'bar',
				markPoint:{
					data:[]
				},
				data:[]
			};
			option.series.push(ser);
			config_cnt++;
			if (config_cnt > 5) {
				option.legend.data.push('');
				config_cnt = 0;
			}
		}

		option.title.text = config_data.title + '-' + config_data.process_name;

		var x_values = [];
		var x_i = 0;
		var testfiles = [];
		for (var i = 1; i < lines.length; i+=step_len) {
			if (lines[i].length < 10) {
				continue;
			}

			var step_time_max = 0;
			var step_mem_max = 0;
			for (var step_i = 0; step_i < step_len; step_i++) {
				var step_str = lines[i + step_i];
				try {
					var data = JSON.parse(step_str);
				} catch (e) {
					console.info(e.message);
					continue;
				}

				// Time
				var time_i = time_config.indexOf(data.info);
				if (time_i != -1) {
					var time_data = data.used_time;
					option.series[time_i].data.push(time_data);
					// console.info(time_i);
					step_time_max += time_data;
					if (step_time_max > time_max) {
						time_max = step_time_max;
					}
					if (time_data > 0) {
						option.series[time_i].markPoint.data.push({name: 'time', value: 100, xAxis: data.testfile, yAxis: time_data});
					}
				}

				// CPU
				var cpu_i = cpu_config.indexOf(data.info)
				if (cpu_i != -1) {
					var cpu_data = data.performance.process[0].cpu_percent;
					var ser_i = cpu_i+time_config.length;

					option.series[ser_i].data.push(cpu_data);
					if (cpu_data > 0) {
						option.series[ser_i].markPoint.data.push({name: 'cpu', value: 100, xAxis: data.testfile, yAxis: cpu_data});
					}
				}

				// Memory
				var mem_i = mem_config.indexOf(data.info)
				if (mem_i != -1) {
					var mem = (data.performance.process[0].memory[0]) / 1024 / 1024;
					mem = Math.round(mem);
					var ser_i = mem_i+cpu_config.length+time_config.length;
					option.series[ser_i].data.push(mem);
					step_mem_max += mem;
					if (step_mem_max > mem_max) {
						mem_max = step_mem_max;
					}
					if (mem > 0) {
						option.series[ser_i].markPoint.data.push({name: 'mem', value: 100, xAxis: data.testfile, yAxis: mem});
					}
				}
				if (x_asix_name == "testfile") {
					if (x_values.indexOf(data.testfile) == -1) {
						x_values.push(data.testfile);
					}
				}


			}
			if (x_asix_name == "index") {
				x_values.push(x_i++);
			}
		}
		option.xAxis[0].data = x_values;
		option.yAxis[0].max = time_max + 100;

		// main_chart.hideLoading();
		// main_chart.setOption(option);
        console.info(option);
        this.logChartOptions.push({name: logname, option : option});

	}

    doRemove(item: any): void {
        for (let i = 0; i < this.logChartOptions.length; i++) {
            if (this.logChartOptions[i].name == item._file.name) {
                this.logChartOptions.splice(i, 1);
                break;
            }
        }
        console.info('Remove:', item);
        item.remove();
    }
}
