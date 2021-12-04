package app

import (
	"context"
	"fmt"
	"github.com/1makarov/binance-nft-buy/internal/pkg/account"
	"github.com/1makarov/binance-nft-buy/internal/pkg/mysterybox"
	bapi "github.com/1makarov/binance-nft-buy/pkg/binance-api"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func App(account *account.Account, box *mysterybox.Box,close context.CancelFunc,taskCtx context.Context) {
	defer fmt.Scanf("\n")
	var checkbotToken,sitekey,cookie,trace string

	body, err := bapi.MarshalMysteryBoxBuy(box.Box.ID, box.Quantity)
	if err != nil {
		log.Fatalf("error marshal buy box: %s\n", err.Error())
	}

	log.Println("Waiting started successfully")
	//wait(box.Information.StartTime-7)
	//wait(1638534700-7)
	chromedp.ListenTarget(
		taskCtx,
		func(ev interface{}){
			if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
				if ev.Request.URL == "https://www.binance.com/bapi/nft/v1/private/nft/nft-trade/product-onsale"{
					sitekey = ev.Request.Headers["x-nft-checkbot-sitekey"].(string)
					checkbotToken = ev.Request.Headers["x-nft-checkbot-token"].(string)
					trace = ev.Request.Headers["x-trace-id"].(string)
					fmt.Println(len(checkbotToken))
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
	_ = chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.WaitVisible(`//button[text()="Подтвердить"]`),
		chromedp.Click(`//button[text()="Подтвердить"]`, chromedp.NodeVisible),
		chromedp.Sleep(3*time.Second),
	)
	req := account.Auth.NFTMysteryBoxGenerateRequest(body,sitekey,checkbotToken,cookie,trace)
	//wait(box.Information.StartTime)
	//wait(1638534700)

	log.Println("Start buy")

	for j := 0; j < 400; j++ {

	for i := 0; i < 5; i++ {
		go func() {
			if !box.Status {
				resp, err := account.Auth.NFTMysteryBoxBuy(req)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println(string(resp.Body()))
				return
			} else {
				return
			}
		}()

	}
	time.Sleep(50 * time.Millisecond)
}


	time.Sleep(6 * time.Second)
	box.Status = true
	time.Sleep(1 * time.Second)

	fmt.Println("Purchases are completed")
}

func wait(s int64) {
	t := time.Unix(s, 0).UTC().Add(-3 * time.Second).Unix()
	for {
		if time.Now().UTC().Unix() >= t {
			return
		}
	}
}
