package coinMarketCap

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
)

func (c *CoinMarketCap) GetCoinData(url string) {

	var httpClient= &http.Client{Timeout: 10 * time.Second} //default client has no timeout set

	response, err := httpClient.Get(url)
	utility.CheckErr(err)

	if err == nil {
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		utility.CheckErr(err)

		err = json.Unmarshal(body, &c.CoinData)
		utility.CheckErr(err)
	}
	return
}

func (c *CoinMarketCap) GetGlobalMarketData(url string) {
	var httpClient = &http.Client{Timeout: 10 * time.Second} //default client has no timeout set

	response, err := httpClient.Get(url)
	utility.CheckErr(err)

	if err == nil {
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		utility.CheckErr(err)

		err = json.Unmarshal(body, &c.GlobalMarketData)
		utility.CheckErr(err)
	}
	return
}
