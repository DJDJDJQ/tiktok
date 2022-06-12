# tiktok

## 抖音项目服务端

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./simple-demo
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试数据

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

### TODO
1. 无法连接数据库抛panic？
2. 代码规范性 
    - 集中结构体之类的到一个文件里
        -token有效性进行封装
    - 目录树，方便记忆和安排
    - 配置文件
        - service.go publish.go 的path
        - mysql.go 的dsn (可以结合pkg/constants.go)
        - 其他
3. 封面图 ffmpeg
4. 云存储 腾讯云cos
5. 在服务器挂着程序运行，关闭命令行或者调试窗口现在就会终止进程
6. 其他未测试的潜在问题 比如上传了很多视频时候的稳定性(需要流畅网速。。。配好cos先)
    - token过期后的错误捕获（前端缺乏退出登录的接口，无法实现迫使前端用户掉线的功能？）
    - 在自己的页面中，点击关注自己可以返回一个不能关注自己的提示？
    - 文件还未上传就把信息写入数据库了
    - 上传视频后，视频流未更新，无法刷到刚上传的视频

### 代码逻辑
1. token有效期24小时(未测试)

2. 上传视频功能需要在服务器端测试