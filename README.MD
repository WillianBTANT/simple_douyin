clone后注意修改`config/config.go`中的数据库配置，程序第一次运行后会生成`config.json`配置文件，可以在`main`函数中进行修改

~~~shell
├─config // 配置文件和程序
├─controller // 接口的handler
├─middleware // 中间件:认证中间件...
├─model // orm model和db连接
├─public // 公开资源
├─test // 单元测试文件
├─utils // 工具包:jwt...
├─router.go // 接口注册
├─.gitignore // gitignore文件
└─main.go // 工具包

~~~
