## Foxit QA automated testing home page


# 功能模块：
+ 可持续集成： ci
    + 各部门ci在线登入、配置、运行入口。
    + ci 任务 RESTFul 控制接口。
+ 测试文件浏览： files-browser
    + 在线查看文件。
    + 在线批量上传文件和设置对应信息。
+ 测试文件搜索： files-search
    + 支持文件名、类型、体积及 SQL 查询。
+ 图表生成： chart
    + 在线动态生成
    + 上传 log 生成。
+ 测试管理： test/state
    + 机器UI连接、性能、用途、可用状态查看。
    + 机器用途在线修改。
    + 测试进度及结果查看。
    + 测试在线中止。
+ 模糊测试： test/fuzz
    + 模糊结果查看。
    + 屏幕直播查看。

## 实现及依赖：
> * 后端： GIN
> * 前端： Angular
> * [其它](https://github.com/xiaoxiayu/foxitqa.cn/blob/master/frontend/package.json)