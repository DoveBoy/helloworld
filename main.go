package main

import (
	"github.com/Albert-Zhan/httpc"
	"github.com/unknwon/goconfig"
	"github.com/ztino/jd_seckill/cmd"
	"github.com/ztino/jd_seckill/common"
	"github.com/ztino/jd_seckill/log"
	"os"
	"runtime"
)

func init()  {
	//软件目录获取
	dir,err:=os.Getwd()
	if err!=nil {
		common.SoftDir="."
	}else{
		common.SoftDir=dir
	}

	//日志初始化
	if !common.IsDir(common.SoftDir+"/logs/") {
		_ = os.Mkdir(common.SoftDir+"/logs/", 0777)
	}

	//客户端设置初始化
	common.Client=httpc.NewHttpClient()
	common.CookieJar=httpc.NewCookieJar()
	common.Client.SetCookieJar(common.CookieJar)

	//配置文件初始化
	confFile:=common.SoftDir+"/conf.ini"
	if common.Config,err=goconfig.LoadConfigFile(confFile);err!=nil {
		log.Println("配置文件不存在，程序退出")
		os.Exit(0)
	}

	//抢购状态管道
	common.SeckillStatus=make(chan bool)
}

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}