package common

import (
	"github.com/Albert-Zhan/httpc"
	"github.com/unknwon/goconfig"
)

const (
	SoftName          = "jd_seckill"
	Version           = "0.1.7"
	DateTimeFormatStr = "2006-01-02 15:04:05"
	DateFormatStr     = "2006-01-02"
)

var Client *httpc.HttpClient

var CookieJar *httpc.CookieJar

var Config *goconfig.ConfigFile

var SeckillStatus chan bool
