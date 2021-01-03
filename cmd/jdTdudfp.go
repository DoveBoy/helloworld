package cmd

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/ztino/jd_seckill/common"
	"github.com/ztino/jd_seckill/jd_seckill"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

func init() {
	rootCmd.AddCommand(jdTdudfpCmd)
}

var jdTdudfpCmd = &cobra.Command{
	Use:   "jdTdudfp",
	Short: "auto get jd eid and fp",
	Run: startJdTdudfp,
}

func startJdTdudfp(cmd *cobra.Command, args []string)  {
	session:=jd_seckill.NewSession(common.CookieJar)
	err:=session.CheckLoginStatus()
	if err!=nil {
		log.Println("自动获取eid和fp失败，请重新登录")
	}else {
		log.Println("开始自动获取eid和fp，如遇卡住请结束进程，重新启动")
		ctx := context.Background()
		options := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", false),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
			chromedp.UserAgent(common.Config.MustValue("config","default_user_agent","")),
		}
		options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

		c, cc := chromedp.NewExecAllocator(ctx, options...)
		defer cc()

		ctx, cancel := chromedp.NewContext(c)
		ch := addNewTabListener(ctx)
		defer cancel()

		u, _ := url.Parse("http://jd.com")
		cookies := common.CookieJar.Cookies(u)
		err := chromedp.Run(ctx,
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
				chromedp.Navigate("https://jd.com"),
				chromedp.Click(".cate_menu_lk"),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		newCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
		ch = addNewTabListener(newCtx)
		defer cancel()

		err = chromedp.Run(newCtx,
			chromedp.Click(`.goods_item_link`),
		)
		if err != nil {
			log.Fatal(err)
		}

		newCtx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
		defer cancel()

		var res []byte
		err = chromedp.Run(newCtx,
			chromedp.Click(`#InitCartUrl`),
			chromedp.WaitVisible(".btn-addtocart"),
			chromedp.Click(".btn-addtocart"),
			chromedp.WaitVisible(".common-submit-btn"),
			chromedp.Click(".common-submit-btn"),
			chromedp.WaitVisible("#sumPayPriceId"),
			chromedp.Sleep(2*time.Second),
			chromedp.Evaluate("_JdTdudfp", &res),
		)
		if err != nil {
			log.Fatal(err)
		}
		value:=string(res)
		if !gjson.Valid(value) {
			log.Println("获取失败，请重新尝试")
		}else{
			log.Println("获取成功，请手动填入配置文件")
			log.Println("eid:"+gjson.Get(value,"eid").String())
			log.Println("fp:"+gjson.Get(value,"fp").String())
		}
	}
}

func addNewTabListener(ctx context.Context) <-chan target.ID {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	return chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return info.URL != ""
	})
}