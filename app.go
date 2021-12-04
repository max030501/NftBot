package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	//"fmt"
	"github.com/1makarov/binance-nft-buy/internal/app"
	acc "github.com/1makarov/binance-nft-buy/internal/domain/account"
	"github.com/1makarov/binance-nft-buy/internal/pkg/account"
	"github.com/1makarov/binance-nft-buy/internal/pkg/mysterybox"
	//"github.com/joho/godotenv"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	//"time"
)

var cookie,csrftoken string
var closeTask context.CancelFunc
var taskCtx context.Context

func main() {

	if err := initEnv(); err != nil {
		log.Println(err)
		closeTask()
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








	app.App(a, box,closeTask,taskCtx)
}

func initEnv() error {
	_ = godotenv.Load()

	//profile_id := "ceb3eba9-5420-47cb-a215-3803c449d49a"
	//url := "http://127.0.0.1:35000/automation/launch/python/"+profile_id


	//opts := append(chromedp.DefaultExecAllocatorOptions[0:2],
	//	chromedp.DefaultExecAllocatorOptions[3:]...)

	//allocCtx, close := chromedp.NewExecAllocator(context.Background(), opts...)
	devtoolsWsURL := flag.String("devtools-ws-url", "ws://127.0.0.1:"+os.Getenv("port"), "DevTools WebSsocket URL")
	flag.Parse()
	if *devtoolsWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	// create allocator context for use with creating a browser context later
	allocCtx, close := chromedp.NewRemoteAllocator(context.Background(), *devtoolsWsURL)
	taskCtx, _ = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	closeTask = close
	chromedp.ListenTarget(
		taskCtx,
		func(ev interface{}){
			if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
				if ev.Request.URL == "https://www.binance.com/bapi/accounts/v1/public/authcenter/auth"{
					csrftoken, _ = ev.Request.Headers["csrftoken"].(string)
					fmt.Println(csrftoken)
				}
			}
			if ev, ok := ev.(*network.EventRequestWillBeSentExtraInfo); ok {
				if val, ok:=ev.Headers["cookie"];ok{
					cookie ,_ =val.(string)
					//fmt.Println(cookie)
				}
			}

		},
	)
	err := chromedp.Run(taskCtx,
		network.Enable(),
		//chromedp.Navigate(`https://accounts.binance.com/ru/login?return_to=aHR0cHM6Ly93d3cuYmluYW5jZS5jb20vcnUvbmZ0L2hvbWU%3D`),
		//chromedp.WaitVisible("//header/div[3]/div[1]"),
		//chromedp.WaitVisible("//button[text()=\"Принять\"]"),
		//chromedp.Click("//button[text()=\"Принять\"]",chromedp.NodeVisible),
		chromedp.Navigate(`https://www.binance.com/ru/nft/goods/sale/161410219850160608?isBlindBox=1&isOpen=false`),
		//chromedp.Navigate("https://www.binance.com/ru/nft/mystery-box/market?page=1&size=16&keyword=&nftType=null&orderBy=amount_sort&orderType=1&serialNo=null&tradeType=1"),
		//chromedp.Sleep(5*time.Second),
		//chromedp.WaitVisible("/html[1]/body[1]/div[1]/div[1]/div[2]/main[1]/div[1]/div[1]/div[4]/div[1]/div[1]/div[1]/div[1]"),
		//chromedp.Click("/html[1]/body[1]/div[1]/div[1]/div[2]/main[1]/div[1]/div[1]/div[4]/div[1]/div[1]/div[1]/div[1]",chromedp.NodeVisible),
		//chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible("//body/div[@id='__APP']/div[1]/div[2]/main[1]/div[1]/div[1]/div[5]/div[2]/div[1]/div[1]/input[1]"),
		chromedp.SendKeys("//body/div[@id='__APP']/div[1]/div[2]/main[1]/div[1]/div[1]/div[5]/div[2]/div[1]/div[1]/input[1]","5"),
		chromedp.WaitVisible(`//button[contains(text(),'Отправить')]`),
		chromedp.Click(`//button[contains(text(),'Отправить')]`, chromedp.NodeVisible),
		)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')

	//err = chromedp.Run(taskCtx,
	//	network.Enable(),
	//	chromedp.Navigate(`https://www.binance.com/ru/nft/goods/mystery-box/detail?productId=14156899&isOpen=true&isProduct=1`),
	//	// wait for footer element is visible (ie, page is loaded)
	//	chromedp.WaitVisible(`//button[text()="Сделать ставку"]`),
	//	// find and click "Example" link
	//	chromedp.Click(`//button[text()="Сделать ставку"]`, chromedp.NodeVisible),
	//
	//)
	return err
}
