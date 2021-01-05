Jd_Seckill
=======
> ⚠ 此项目是[python jd_seckill](https://github.com/huanghyw/jd_seckill) 的go版本实现，旨在降低使用门栏和相互学习而创建。

**go版本的jd_seckill，京东抢茅台神器，支持跨平台，使用者请在发布页下载可执行文件，欢迎pr。**

## 支持系统

>目前编译好的可执行文件有Windows,MacOS,Linux,arm,mips平台。

## 安装

方式一(推荐):

```shell
git clone https://github.com/ztino/jd_seckill.git
cd jd_seckill
go get
```

方式二:

```shell
go get github.com/ztino/jd_seckill
```

## 待办

- 跨平台桌面端支持，打算使用：https://github.com/therecipe/qt

## 使用

> [下载](https://github.com/ztino/jd_seckill/releases) 对应平台的可执行文件，解压，终端进入该目录。

> ⚠ 0.1.4(包含)之前版本，不适用该教程，请直接执行命令，按照提示操作。

### 登录
执行以下命令按照提示操作:
```shell
jd_seckill login
```

### 自动获取eid,fp

> ⚠依赖谷歌浏览器，请安装谷歌浏览器，windows下请将安全目录加入系统变量Path
> ⚠【重要提醒】自动获取eid和fp期间，建议鼠标跟随页面跳转，滑动到【加入购物车】【去购车结算】【去结算】按钮，但不要点击，可以提升获取成功率！(经验之谈！)

执行以下命令按照提示操作:
```shell
jd_seckill jdTdudfp
```
> ⚠目前实验性阶段，请勿依赖该功能

### 预约
执行以下命令按照提示操作:
```shell
jd_seckill reserve
```

### 抢购
执行以下命令按照提示操作:
```shell
jd_seckill seckill
```

### 退出登录
```shell
jd_seckill logout
```

### 获取版本号
```shell
jd_seckill version
```

> ⚠ 以上命令并不是每次都需要执行的，都是可选的，具体使用请参考提示。

### Linux下命令行方式显示二维码（以Ubuntu为例）

```bash
$ sudo apt-get install qrencode zbar-tools # 安装二维码解析和生成的工具，用于读取二维码并在命令行输出。
$ zbarimg qr_code.png > qrcode.txt && qrencode -r qrcode.txt -o - -t UTF8 # 解析二维码输出到命令行窗口。
```

## 使用教程

#### 1. 推荐Chrome浏览器
#### 2. 网页扫码登录，或者账号密码登录
#### 3. 填写config.ini配置信息
(1)`eid`和`fp`找个普通商品随便下单,然后抓包就能看到,这两个值可以填固定的
> 随便找一个商品下单，然后进入结算页面，打开浏览器的调试窗口，切换到控制台Tab页，在控制台中输入变量`_JdTdudfp`，即可从输出的Json中获取`eid`和`fp`。  
> 不会的话参考issue https://github.com/ztino/jd_seckill/issues/2

(2)`sku_id`,`default_user_agent`
> `sku_id`已经按照茅台的填好。
> `default_user_agent` 可以用默认的。谷歌浏览器也可以浏览器地址栏中输入about:version 查看`USER_AGENT`替换

(3)配置一下时间
> 现在不强制要求同步最新时间了，程序会自动同步京东时间
>> 但要是电脑时间快慢了好几个小时，最好还是同步一下吧

以上都是必须的.
> tips：
> 在程序开始运行后，会检测本地时间与京东服务器时间，输出的差值为本地时间-京东服务器时间，即-50为本地时间比京东服务器时间慢50ms。
> 本代码的执行的抢购时间以本地电脑/服务器时间为准

(4)修改抢购瓶数
> 可在配置文件中找到seckill_num进行修改

(5)其他配置
> 请自行参考使用

## 抢购流程/抢购结果

- 程序开始抢购总时间为两分钟，不管有无抢购成功，都会停止，抢购详情请查阅日志和自己配置的第三方推送服务(目前支持：钉钉机器人、SMTP邮件、Server酱)

- 抢购开始时间设定，格式：09:59:59，如果小于当前时间，则表示明天这个时间点开始；大于当前时间，则为今天

- 先写这么多。。。

## 感谢
##### 非常感谢原作者 https://github.com/zhou-xiaojun/jd_mask 提供的代码
##### 也非常感谢 https://github.com/wlwwu/jd_maotai 进行的优化
