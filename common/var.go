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

var (
	winExecError = map[uint32]string{
		0:  "The system is out of memory or resources.",
		2:  "The .exe file is invalid.",
		3:  "The specified file was not found.",
		11: "The specified path was not found.",
	}
)