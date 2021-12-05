package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	"github.com/max030501/NftBot/internal/app"
	acc "github.com/max030501/NftBot/internal/domain/account"
	"github.com/max030501/NftBot/internal/pkg/account"
	"github.com/max030501/NftBot/internal/pkg/mysterybox"
	"log"
	"os"
	"strings"
	"time"
)

var cookie, csrftoken string
var closeTask context.CancelFunc
var taskCtx context.Context

func main() {

	if err := initEnv(); err != nil {
		log.Println(err)
		return
	}

	a, err := account.InitAccount(acc.Setting{
		Proxy: os.Getenv("PROXY"),
		BAuth: &acc.BAuth{Cookie: cookie, Csrf: csrftoken},
	})

	if err != nil {
		log.Println(err)
		closeTask()
		return
	}
	if err = a.HandleAccount(); err != nil {
		log.Println(err)
		closeTask()
		return
	}

	boxList, err := mysterybox.GetActiveMysteryBoxList()
	if err != nil {
		log.Println(err)
		return
	}

	box, err := boxList.SelectBox()
	if err != nil {
		log.Println(err)
		return
	}

	if err = box.InitBox(); err != nil {
		log.Println(err)
		return
	}

	app.App(a, box, closeTask, taskCtx)
}

func initEnv() error {
	_ = godotenv.Load()
	var allocCtx context.Context
	if os.Getenv("port") == "" {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)
		allocCtx, closeTask = chromedp.NewExecAllocator(context.Background(), opts...)
	} else {
		allocCtx, closeTask = chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:"+os.Getenv("port"))
	}
	taskCtx, _ = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	var reqID network.RequestID
	auth := false
	chromedp.ListenTarget(
		taskCtx,
		func(ev interface{}) {
			if ev, ok := ev.(*network.EventLoadingFinished); ok {
				if reqID == ev.RequestID {
					go func() {
						c := chromedp.FromContext(taskCtx)
						rbp := network.GetResponseBody(ev.RequestID)
						body, err := rbp.Do(cdp.WithExecutor(taskCtx, c.Target))
						if err == nil && strings.Contains(string(body), "{\"code\":\"000000\",\"data\":60,\"success\":true}") {
							auth = true
						}
					}()
				}
			}
			if ev, ok := ev.(*network.EventResponseReceived); ok {
				if ev.Response.URL == "https://www.binance.com/bapi/accounts/v1/public/authcenter/auth" {
					reqID = ev.RequestID

				}
			}
			if ev, ok := ev.(*network.EventRequestWillBeSentExtraInfo); ok {
				if csrf, ok := ev.Headers["csrftoken"]; ok && csrftoken == "" && auth {
					cookie, _ = ev.Headers["cookie"].(string)
					csrftoken = csrf.(string)
				}
			}
		},
	)
	err := chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(`https://www.binance.com/ru/nft/home`),
		chromedp.Sleep(time.Second),
	)
	if !auth {
		err = chromedp.Run(taskCtx,
			network.Enable(),
			chromedp.Navigate(`https://accounts.binance.com/ru/login?return_to=aHR0cHM6Ly93d3cuYmluYW5jZS5jb20vcnUvbmZ0L2hvbWU%3D`),
			chromedp.WaitVisible("//a[contains(text(),'Мистери-боксы')]"),
		)
	}
	err = chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(`https://www.binance.com/ru/nft/goods/sale/161410219850160608?isBlindBox=1&isOpen=false`),
		chromedp.WaitVisible("//body/div[@id='__APP']/div[1]/div[2]/main[1]/div[1]/div[1]/div[5]/div[2]/div[1]/div[1]/input[1]"),
		chromedp.SendKeys("//body/div[@id='__APP']/div[1]/div[2]/main[1]/div[1]/div[1]/div[5]/div[2]/div[1]/div[1]/input[1]", "5"),
		chromedp.WaitVisible(`//button[contains(text(),'Отправить')]`),
		chromedp.Click(`//button[contains(text(),'Отправить')]`, chromedp.NodeVisible),
	)

	return err
}
