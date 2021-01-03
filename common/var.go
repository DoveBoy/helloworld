package common

import (
	"github.com/Albert-Zhan/httpc"
	"github.com/unknwon/goconfig"
)

const SoftName = "jd_seckill"

const Version = "0.1.6"

var Client *httpc.HttpClient

var CookieJar *httpc.CookieJar

var Config *goconfig.ConfigFile

var SeckillStatus chan bool
