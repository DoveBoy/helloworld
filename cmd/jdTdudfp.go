package cmd

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/unknwon/goconfig"
	"github.com/ztino/jd_seckill/common"
	"github.com/ztino/jd_seckill/jd_seckill"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

func init() {
	rootCmd.AddCommand(jdTdudfpCmd)
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
		retryTimes, _ := strconv.Atoi(common.Config.MustValue("config", "retry_times", "5"))

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
		_ = addNewTabListener(ctx)
		defer cancel()

		//设置cookie
		u, _ := url.Parse("http://jd.com")
		cookies := common.CookieJar.Cookies(u)
		err = chromedp.Run(ctx,
			chromedp.Tasks{
				chromedp.ActionFunc(func(ctx context.Context) error {
					expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
					for _, cookie := range cookies {
						_, _ = network.SetCookie(cookie.Name, cookie.Value).
							WithExpires(&expr).
							WithPath("/").
							WithDomain("." + cookie.Domain).
							Do(ctx)
					}
					return nil
				}),
			},
		)
		if err != nil {
			log.Println("设置cookie出错了")
			log.Fatal(err)
		}

	RETRY:
		retryTimes--
		log.Println("【重要提醒】自动获取eid和fp期间，建议鼠标跟随页面跳转，滑动到【加入购物车】【去购车结算】【去结算】按钮，但不要点击，可以提升获取成功率！")
		log.Println(fmt.Sprintf("开始自动获取eid和fp，如遇卡住请耐心等待，重试次数剩余： %v 次", retryTimes))
		var res []byte
		testSkuId := common.Config.MustValue("config", "test_sku_id", "")
		err = chromedp.Run(ctx,
			chromedp.Tasks{
				chromedp.Navigate(fmt.Sprintf("http://item.jd.com/%s.html", testSkuId)),
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
			},
		)
		if err != nil {
			log.Println("chromedp 出错了")
			log.Fatal(err)
		}

		value := string(res)
		if !gjson.Valid(value) || gjson.Get(value, "eid").String() == "" || gjson.Get(value, "fp").String() == "" {
			log.Println("获取失败，请重新尝试，返回信息:" + value)
			if retryTimes > 0 {
				goto RETRY
			}
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

func addNewTabListener(ctx context.Context) <-chan target.ID {
	return chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return true
	})
}
