Jd_Seckill
=======

#### 还未收到安全通知，但是估计也快了，给大家建个Telegram群吧，茶余饭后，聊天扯淡，有条件的可以加入：https://t.me/joinchat/GsDnhtkdKJ4nbwJh

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
- 自动化预约抢购支持，程序自动去茅台页面获取下一次抢购时间
- 跨平台桌面端支持，打算使用：https://github.com/therecipe/qt

## 使用

> [下载](https://github.com/ztino/jd_seckill/releases) 对应平台的可执行文件，解压，终端进入该目录。

### 登录
执行以下命令按照提示操作:
```shell
jd_seckill login
```

### 自动获取eid,fp

> ⚠依赖谷歌浏览器，请安装谷歌浏览器，windows下请将安装目录加入系统变量Path

> ⚠ 京东可能在修改eid和fp的获取方式了，目前该功能获取不太稳定，请勿依赖该功能，目前观望中，不做更改

执行以下命令按照提示操作:
```shell
#参数--good_url商品链接必须设置，链接地址是一个可以加入购物车的商品
jd_seckill jdTdudfp --good_url https://item.jd.com/100007959916.html
```
> ⚠获取成功后会将获取到的eid和fp写入到配置文件中

### 预约
执行以下命令按照提示操作:
```shell
jd_seckill reserve
```

### 抢购
执行以下命令按照提示操作:
```shell
#支持--run参数，将跳过抢购等待时间，直接执行抢购任务，适合10点左右未设置抢购时间的使用
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

## 使用教程

#### 1. 推荐Chrome浏览器
#### 2. 网页扫码登录，或者账号密码登录
#### 3. 填写config.ini配置信息

> ⚠ 按照下方获取不到的，可以点击进入付款界面(输入支付密码页面)，尝试下方步骤进行获取

(1)`eid`和`fp`找个普通商品随便下单,然后抓包就能看到,这两个值可以填固定的
> 随便找一个商品下单，然后进入结算页面，打开浏览器的调试窗口，切换到控制台Tab页，在控制台中输入变量`_JdTdudfp`，即可从输出的Json中获取`eid`和`fp`。  
> 不会的话参考issue https://github.com/ztino/jd_seckill/issues/2

(2)`sku_id`,`default_user_agent`
> `sku_id`已经按照茅台的填好。
> `default_user_agent` 可以用默认的。谷歌浏览器也可以浏览器地址栏中输入about:version 查看`USER_AGENT`替换

(3)配置一下时间
> 现在不强制要求同步最新时间了，程序会自动同步京东时间
> 但要是电脑时间快慢了好几分钟的，最好还是同步一下吧

以上都是必须的.
> tips：
> 在程序开始运行后，会检测本地时间与京东服务器时间，输出的差值为本地时间-京东服务器时间，即-50为本地时间比京东服务器时间慢50ms。
> 本代码的执行的抢购时间以本地电脑/服务器时间为准

> ⚠ 京东每月限购两瓶，如果本月已抢到两瓶，一个月后再抢吧，有的抢到1瓶的，使用脚本记得需要修改参数

(4)修改抢购瓶数
> 可在配置文件中找到seckill_num进行修改，默认值2瓶

(5)抢购总时间
> 可在配置文件中找到seckill_time进行修改，单位:分钟，默认两分钟

(6)抢购任务数量
> 可在配置文件中找到task_num进行修改，默认5个

(7)每次抢购间隔时间
> 可在配置文件中找到ticker_time进行修改，单位:毫秒，默认1500毫秒，每1000毫秒等于1秒

(8)通知配置
> 目前支持email，wechat，dingtalk，具体可查看配置文件

## 感谢
##### 非常感谢原作者 https://github.com/zhou-xiaojun/jd_mask 提供的代码
##### 也非常感谢 https://github.com/wlwwu/jd_maotai 进行的优化
