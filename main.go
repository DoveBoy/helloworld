package main

import (
	"errors"
	"fmt"
	"github.com/Albert-Zhan/httpc"
	"github.com/tidwall/gjson"
	"github.com/unknwon/goconfig"
	"github.com/ztino/jd_seckill/jd_seckill"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var client *httpc.HttpClient

var cookieJar *httpc.CookieJar

var config *goconfig.ConfigFile

var seckillStatus chan bool

func init()  {
	//客户端设置初始化
	client=httpc.NewHttpClient()
	cookieJar=httpc.NewCookieJar()
	client.SetCookieJar(cookieJar)
	client.SetRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})

	//配置文件初始化
	confFile:="./conf.ini"
	var err error
	if config,err=goconfig.LoadConfigFile(confFile);err!=nil {
		log.Println("配置文件不存在，程序退出")
		os.Exit(0)
	}

	//抢购状态管道
	seckillStatus=make(chan bool)
}

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//用户登录
	user:=jd_seckill.NewUser(client,config)
	wlfstkSmdl,err:=user.QrLogin()
	if err!=nil{
		os.Exit(0)
	}
	ticket:=""
	for  {
		ticket,err=user.QrcodeTicket(wlfstkSmdl)
		if err==nil && ticket!=""{
			break
		}
		time.Sleep(2*time.Second)
	}
	_,err=user.TicketInfo(ticket)
	if err==nil {
		log.Println("登录成功")
		//刷新用户状态
		if status:=user.RefreshStatus();status==nil {
			//活跃用户会话,当会话失效自动退出程序
			go KeepSession(user)
			//获取用户信息
			userInfo,_:=user.GetUserInfo()
			log.Println("用户:"+userInfo)
			//开始预约,预约过的就重复预约
			seckill:=jd_seckill.NewSeckill(client,config)
			seckill.MakeReserve()
			//计算抢购时间
			nowLocalTime:=time.Now().UnixNano()/1e6
			jdTime,_:=GetJdTime()
			buyDate:=config.MustValue("config","buy_time","")
			loc, _ := time.LoadLocation("Local")
			t,_:=time.ParseInLocation("2006-01-02 15:04:05",buyDate,loc)
			buyTime:=t.UnixNano()/1e6
			diffTime:=nowLocalTime-jdTime
			log.Println(fmt.Sprintf("正在等待到达设定时间:%s，检测本地时间与京东服务器时间误差为【%d】毫秒",buyDate,diffTime))
			timerTime:=(buyTime+diffTime)-jdTime
			if timerTime<=0 {
				log.Println("请设置抢购时间")
				os.Exit(0)
			}
			//等待抢购
			time.Sleep(time.Duration(timerTime)*time.Millisecond)
			//开始抢购
			log.Println("时间到达，开始执行……")
			//开启抢购任务,第二个参数为开启几个协程
			//怕封号的可以减少协程数量,相反抢到的成功率也减低了
			Start(seckill,5)
		}else{
			log.Println("登录失效")
		}
	}else{
		log.Println("登录失败")
	}
}

func GetJdTime() (int64,error) {
	req:=httpc.NewRequest(client)
	resp,body,err:=req.SetUrl("https://a.jd.com//ajax/queryServerData.html").SetMethod("get").Send().End()
	if err!=nil || resp.StatusCode!=http.StatusOK {
		log.Println("获取京东服务器时间失败")
		return 0,errors.New("获取京东服务器时间失败")
	}
	return gjson.Get(body,"serverTime").Int(),nil
}

func Start(seckill *jd_seckill.Seckill,taskNum int)  {
	seckillTotalTime:=time.Now().Add(2*time.Minute).Unix()
	//开始检测抢购状态
	go CheckSeckillStatus()
	//抢购总时间两分钟,超时程序自动退出
	for time.Now().Unix()<seckillTotalTime {
		for i:=1;i<=taskNum;i++ {
			go task(seckill)
		}
		//每隔1.5秒执行一次
		//怕封号的可以增加间隔时间,相反抢到的成功率也减低了
		time.Sleep(1500*time.Millisecond)
	}
	log.Println("抢购结束，具体详情请查看日志")
}

func task(seckill *jd_seckill.Seckill)  {
	seckill.RequestSeckillUrl()
	seckill.SeckillPage()
	flag:=seckill.SubmitSeckillOrder()
	//提前抢购成功的,直接结束程序
	if flag {
		//通知管道
		seckillStatus<-true
	}
}

func CheckSeckillStatus()  {
	for {
		select {
			case <-seckillStatus:
			//抢购成功,程序退出
			os.Exit(0)
		}
	}
}

func KeepSession(user *jd_seckill.User)  {
	//每30分钟检测一次
	t:=time.NewTicker(30*time.Minute)
	for {
		select {
		case <-t.C:
			if err:=user.RefreshStatus();err!=nil {
				log.Println("会话失效,程序自动退出")
				os.Exit(0)
			}
			log.Println("活跃会话成功")
		}
	}
}