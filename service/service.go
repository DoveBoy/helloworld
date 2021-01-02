package service

import "github.com/unknwon/goconfig"

func SendMessage(conf *goconfig.ConfigFile,title,msg string) error {
	if conf.MustValue("messenger","enable","false")=="true" {
		//邮件发送
		if conf.MustValue("messenger","type","none")=="smtp" {
			email:=NewEmail(conf)
			err:=email.Send([]string{conf.MustValue("messenger","email","")},title,msg)
			return err
		}
		//Server酱推送
		if conf.MustValue("messenger","type","none")=="wechat" {
			wechat:=NewWechat(conf)
			err:=wechat.Send(title,msg)
			return err
		}
	}
	return nil
}