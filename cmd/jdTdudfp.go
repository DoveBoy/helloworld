package cmd

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/unknwon/goconfig"
	"github.com/ztino/jd_seckill/common"
	"github.com/ztino/jd_seckill/jd_seckill"
	"github.com/ztino/jd_seckill/log"
	"net/url"
	"os"
	"time"
)

func init() {
	rootCmd.AddCommand(jdTdudfpCmd)
	jdTdudfpCmd.Flags().StringP("good_url","g","","")
	_=jdTdudfpCmd.MarkFlagRequired("good_url")
}

var jdTdudfpCmd = &cobra.Command{
	Use:   "jdTdudfp",
	Short: "auto get jd eid and fp",
	Run:   startJdTdudfp,
}

func startJdTdudfp(cmd *cobra.Command, args []string) {
	session := jd_seckill.NewSession(common.CookieJar)
	err := session.CheckLoginStatus()
	if err != nil {
		log.Println("自动获取eid和fp失败，请重新登录")
	} else {
		log.Println("开始自动获取eid和fp，如遇卡住请结束进程，重新启动")
		options := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", false),                       //debug使用
			chromedp.Flag("blink-settings", "imagesEnabled=false"), //禁用图片加载
			chromedp.Flag("start-maximized", true),                 //最大化窗口
			chromedp.Flag("no-sandbox", true),                      //禁用沙盒, 性能优先
			chromedp.Flag("disable-setuid-sandbox", true),          //禁用setuid沙盒, 性能优先
			chromedp.Flag("no-default-browser-check", true),        //不检查默认浏览器
			chromedp.Flag("disable-plugins", true),                 //禁用扩展
			chromedp.UserAgent(common.Config.MustValue("config", "default_user_agent", "")),
		}
		options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

		c, cc := chromedp.NewExecAllocator(context.Background(), options...)
		defer cc()

		ctx, cancel := chromedp.NewContext(c)
		defer cancel()

		//获取cookie
		u, _ := url.Parse("http://jd.com")
		cookies := common.CookieJar.Cookies(u)

		//商品链接
		good_url,_:=cmd.Flags().GetString("good_url")

		var res []byte
		err = chromedp.Run(ctx,
			chromedp.Tasks{
				chromedp.ActionFunc(func(ctx context.Context) error {
					//设置cookie
					expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
					for _, cookie := range cookies {
						network.SetCookie(cookie.Name, cookie.Value).
							WithExpires(&expr).
							WithPath("/").
							WithDomain("." + cookie.Domain).
							Do(ctx)
					}
					return nil
				}),
			},
			chromedp.Navigate(good_url),
			chromedp.WaitVisible("#InitCartUrl"), //加入购物车
			chromedp.Sleep(2 * time.Second),
			chromedp.Click("#InitCartUrl"),
			chromedp.WaitVisible(".btn-addtocart"), //去购车结算
			chromedp.Sleep(2 * time.Second),
			chromedp.Click(".btn-addtocart"),
			chromedp.WaitVisible(".common-submit-btn"), //去结算
			chromedp.Sleep(2 * time.Second),
			chromedp.Click(".common-submit-btn"),
			chromedp.Sleep(3 * time.Second),
			chromedp.Evaluate("_JdTdudfp", &res),
		)
		if err != nil {
			log.Println("chromedp 出错了")
			log.Fatal(err)
		}

		value := string(res)
		if !gjson.Valid(value) || gjson.Get(value, "eid").String() == "" || gjson.Get(value, "fp").String() == "" {
			log.Println("获取失败，请重新尝试，返回信息:" + value)
		} else {
			eid := gjson.Get(value, "eid").String()
			fp := gjson.Get(value, "fp").String()
			log.Println("eid:" + eid)
			log.Println("fp:" + fp)

			//修改配置文件
			confFile := "./conf.ini"
			cfg, err := goconfig.LoadConfigFile(confFile)
			if err != nil {
				log.Println("配置文件不存在，程序退出")
				os.Exit(0)
			}

			cfg.SetValue("config", "eid", eid)
			cfg.SetValue("config", "fp", fp)
			if err := goconfig.SaveConfigFile(cfg, confFile); err != nil {
				log.Println("保存配置文件失败，请手动填入配置文件")
			}

			log.Println("eid, fp参数已经自动填入配置文件")
		}

	}
}
