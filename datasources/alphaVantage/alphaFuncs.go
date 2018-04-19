package alphaVantage

import (
	"net/http"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"io/ioutil"
	"encoding/json"
	"strings"
)

func GetTop20ETFs(url string, etf string) [][]string {
	var httpClient = &http.Client{Timeout: 10 * time.Second} //default client has no timeout set
	response, err := httpClient.Get(url)
	utility.CheckErr(err)

	result := [][]string {nil}
	if err == nil {
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		utility.CheckErr(err)

		jsonRaw := map[string]json.RawMessage{}
		err = json.Unmarshal(body, &jsonRaw)
		utility.CheckErr(err)

		/*
		resp := Response{} // do we need this section of codefrom HERE...

		if err := json.Unmarshal(body, &resp); err != nil {
			utility.LogToFile(err.Error())
			return nil
		}					//...to HERE? It doesnt do anything after all
		*/
		var jsonMap jsonMap

		if err := json.Unmarshal(jsonRaw["Time Series (Daily)"], &jsonMap); err != nil {
			utility.LogToFile(err.Error())
			return nil
		}

		for elem := range jsonMap {
			result = append(result, []string {
				elem,
				strings.Split(etf, "=")[1],
				jsonMap[elem].Open,
				jsonMap[elem].Close,
				jsonMap[elem].High,
				jsonMap[elem].Low,
				jsonMap[elem].Volume,
				})
		}
	}
	return result
}
