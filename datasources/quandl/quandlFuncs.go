package quandl

import (
	"net/http"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"io/ioutil"
	"encoding/json"
)

func (gd *Quandl) GetLMBAGold(url string) GoldDataTimeSeries {
	var httpClient = &http.Client{Timeout: 10 * time.Second} //default client has no timeout set
	response, err := httpClient.Get(url)
	utility.CheckErr(err)

	if err == nil {
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		utility.CheckErr(err)

		qr := Response{}
		err = json.Unmarshal(body, &qr)
		utility.CheckErr(err)

		dataArray := qr.Data.([]interface{})

		for idx := range dataArray {
			switch vv := dataArray[idx].(type) {
				case []interface{}: {
					switch vv[0].(type) {
						case string: {
							for idx2 := range vv{
								if vv[idx2] == nil {
									vv[idx2] = 0.0 //we dont need to care about the case where reflect.TypeOf(vv[idx2]) returns string because strings cant be nil anyways
								}
							}
							gd.GoldDataTimeSeries.GoldDataEntries = append(gd.GoldDataTimeSeries.GoldDataEntries, GoldDataEntry{
								vv[0].(string),
								vv[1].(float64),
								vv[2].(float64),
								vv[3].(float64),
								vv[4].(float64),
								vv[5].(float64),
								vv[6].(float64),
							})
							//fmt.Println(gd.GoldDataTimeSeries.GoldDataEntries)
						}
						default: // do something
					}
				}
				default: //do something
			}
		}
	}
	return gd.GoldDataTimeSeries
}
