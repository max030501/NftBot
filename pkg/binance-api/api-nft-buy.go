package bapi

import (
	"encoding/json"
	binance_struct "github.com/max030501/NftBot/internal/domain/binance-api"
	"github.com/valyala/fasthttp"
)

const (
	urlNFTMysteryBoxBuy = "https://www.binance.com/bapi/nft/v1/private/nft/mystery-box/purchase"
)

func (api *Api) NFTMysteryBoxGenerateRequest(body []byte, sitekey string, checkbot string) *fasthttp.Request {
	r := fasthttp.AcquireRequest()
	api.request.CopyTo(r)
	r.Header.SetMethod(fasthttp.MethodPost)
	r.Header.SetContentType("application/json")
	r.Header.SetRequestURI(urlNFTMysteryBoxBuy)
	r.SetBody(body)
	r.Header.Set("x-nft-checkbot-sitekey", sitekey)
	r.Header.Set("x-nft-checkbot-token", checkbot)
	return r
}

func (api *Api) NFTMysteryBoxBuy(req *fasthttp.Request) (*fasthttp.Response, error) {
	response, err := api.postRequest(req)
	if err = handleError(response, err); err != nil {
		return nil, err
	}
	return response, nil
}

func MarshalMysteryBoxBuy(productID string, amount int) ([]byte, error) {
	b, err := json.Marshal(binance_struct.BuyRequest{ProductID: productID, Amount: amount})
	if err != nil {
		return nil, err
	}
	return b, nil
}
