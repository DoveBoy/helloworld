package main

import (
	"github.com/Albert-Zhan/httpc"
	"github.com/unknwon/goconfig"
	"github.com/ztino/jd_seckill/cmd"
	"github.com/ztino/jd_seckill/common"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

func init()  {
	//将日志同时输出到控制台和文件
	file := "./" + "jd_seckill_" + time.Now().Format("20060102") + ".log"
	logFile, logErr := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if logErr != nil {
		panic(logErr)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetPrefix("[jd_seckill]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	//客户端设置初始化
	common.Client=httpc.NewHttpClient()
	common.CookieJar=httpc.NewCookieJar()
	common.Client.SetCookieJar(common.CookieJar)

	//配置文件初始化
	confFile:="./conf.ini"
	var err error
	if common.Config,err=goconfig.LoadConfigFile(confFile);err!=nil {
		log.Println("配置文件不存在，程序退出")
		os.Exit(0)
	}

	//抢购状态管道
	common.SeckillStatus=make(chan bool)
	log.Println("jd_seckill 程序启动成功，祝您成功！")
}

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}