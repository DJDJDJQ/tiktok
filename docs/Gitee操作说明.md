# Git操作说明

## 0.规范和建议

### git常用结构和规范

![image-20220611115356872](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611115356872.png)

commit的时候 采用 <type>: <subject>的格式 <subject>中写具体的修改内容

如 fix: 无法取消点赞

##### type类型

1. feat         增加新特性

2. fix            修复bug

3. docs        修改文档（readme.md）

4. style        修改格式不改变代码逻辑

5. refactor  代码重构，没有新功能或bug

6. perf         性能测试

7. test          测试用例

8. chore      增加依赖库 工具等

### 建议

- 放弃更改之后，代码会回滚到更改之前的状态，如果在其中修改了部分代码会直接丢失到上一个状态（深痛教训）
- 如果出错可以通过ref回滚代码（仅限已经commit的历史版本）
- 在上传代码前请仔细检查没有覆盖别人的代码

## 1.Gitee上的操作

首先注册Gitee并加入组织

使用邀请链接

https://gitee.com/organizations/lof-bytedance/invite?invite=3d3439873b35f1cdb26538a98495c381bf30ee404b2461e175d29b6a21b2efb7d03274eeb9ac7ad163ce7e16acfba504

需要经过同意后进入组织，就可以看到仓库了

## 2.git的相关配置

*下面的部分参考了这篇博客 https://blog.csdn.net/weixin_46571373/article/details/107525877*

*Goland环境可以参考这篇博客 https://blog.csdn.net/qq_42956653/article/details/121613703*

#### 打开git bash，设置全集变量为自己注册码云的用户名和邮箱

```shell
git config --global user.name yourname # "你码云的名字或昵称"
git config --global user.email youremail@xxx.com # "你码云的主邮箱"
```

#### 添加之后，在gitee中添加自己的ssh链接

``` bash
ssh-keygen -t rsa -C youremail@xxx.com
cat ~/.ssh/id_rsa.pub
```

将秘钥粘贴到gitee中

![image-20220611114723483](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114723483.png)![image-20220611114739781](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114739781.png)

## ide中的配置(vscode)

打开一个空白vscode页面，在左边分支管理中选择克隆存储库

输入链接https://gitee.com/lof-bytedance/tiktok 

选择一个位置存放pull下来的项目文件

![image-20220611114850915](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114850915.png)

会要求登录gitee账户

![image-20220611114900721](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114900721.png)

这样就将文件pull到本地了

然后可以将自己的文件全部覆盖到down下来的文件夹中，逐一进行比对，

是自己的内容就进行暂存，不是自己的内容就放弃更改

![image-20220611114910772](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114910772.png)

![image-20220611114916039](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114916039.png)

点击对勾进行push到自己本地的仓库中

![image-20220611114926170](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114926170.png)

点击同步更改后就完成了一次push

![image-20220611114934817](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114934817.png)

在进行了修改以后，可以从master中down下来所有进行的修改

![image-20220611114946504](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114946504.png)

![image-20220611114949818](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114949818.png)

![image-20220611114954169](C:\Users\DJDJDJQ\AppData\Roaming\Typora\typora-user-images\image-20220611114954169.png)

