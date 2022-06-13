# 腾讯云COS对象存储
### Q&A
1. 权限问题 私有读写、公有读私有写和公有读写
    - 想法：
        - 在读写时修改权限 ×可能有并发问题和锁的问题
        - 直接修改成公有读写 !危险 暴露url可能存在安全问题，产生意外流量
        - 将存储视频的文件夹权限设置为公有读私有写了 !仍存在流量风险
2. api环境变量

    export SECRETID_TC=AKIDOfjY9WfZEepkPn0yh9MK2Uiv12yunY5N

    export SECRETKEY_TC=81f4FEyTknbBWbR1LQFpAG9mBQSGqJXV

    由于设置了公有读私有写，只能在服务器运行程序时才能上传视频，浏览视频不影响
    - 在获取密钥时可以权限管理 获取临时密钥而不是将永久密钥提前写道环境变量？
3. 存储桶查询可以做一个查询？从而可以更改baseurl

### 踩坑
1. linux环境变量 export语句仅对当前shell有效 关闭shell失效
    - 永久有效需修改~/.bashrc
    - 未解决

### 封面获取

- 通过腾讯云数据万象服务，创建一个视频截帧工作流
- 截图命名为
        ```
        ${InputName}_${Number}.jpg
        ```
- 默认的视频截帧模板不适合获取单张封面，创建一个视频封面ShootFirst
- 配置
    - 设置0秒时截图一张，默认截取每帧，单视频只截图一张，最后输出原图
    - 增加了黑屏判断(未测试)
        - 开启后将会检测视频的前五秒是否存在黑屏，若按指定截帧方式会截取到前5秒内的帧，则从前五秒中第一个非黑屏的帧开始截取。
        - 当黑色像素占比90以上，判断为黑屏，不进行截图
        - 若全程黑屏会怎样？？
- 在完成视频流以后，COS将会回调一个函数到服务器，配置router以后捕获参数获取存储封面图的路径
- 现在是暴力凭借url获取封面链接。。。

### 参考链接
- [腾讯云访问控制](https://cloud.tencent.com/document/product/436/30749)

- [GO SDK](https://cloud.tencent.com/document/product/436/31215)

- [权限问题](https://cloud.tencent.com/document/product/436/13312#.E6.9D.83.E9.99.90.E7.B1.BB.E5.88.AB)

- [密钥管理](https://console.cloud.tencent.com/cam/capi)

- [上传文件](https://blog.csdn.net/qq_62345961/article/details/124480151)

- [截取视频封面](https://cloud.tencent.com/document/product/460/47505)

- [数据万象工作流变量名称](https://cloud.tencent.com/document/product/436/53967#1)