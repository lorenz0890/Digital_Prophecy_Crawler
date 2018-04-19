package poloniex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
)

func (p *Poloniex)GetChartData(url string) {

	var httpClient= &http.Client{Timeout: 10 * time.Second} //default client has no timeout set

	response, err := httpClient.Get(url)
	utility.CheckErr(err)
	if err == nil {

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		utility.CheckErr(err)

		err = json.Unmarshal(body, &p.ChartData)
		utility.CheckErr(err)
	}

	return
}
