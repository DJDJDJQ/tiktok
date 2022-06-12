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
    - 在获取密钥时可以权限管理 获取临时密钥而不是将永久密钥提前写道环境变量？
3. 存储桶查询可以做一个查询？从而可以更改baseurl

4. 没看懂service中是干啥的

### 踩坑
1. linux环境变量 export语句仅对当前shell有效 关闭shell失效
    - 永久有效需修改~/.bashrc
    - 未解决

### 参考链接
[腾讯云访问控制](https://cloud.tencent.com/document/product/436/30749)

[GO SDK](https://cloud.tencent.com/document/product/436/31215)

[权限问题](https://cloud.tencent.com/document/product/436/13312#.E6.9D.83.E9.99.90.E7.B1.BB.E5.88.AB)

[密钥管理](https://console.cloud.tencent.com/cam/capi)

[上传文件](https://blog.csdn.net/qq_62345961/article/details/124480151)