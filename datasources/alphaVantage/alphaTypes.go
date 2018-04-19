package alphaVantage

type MetaData struct {
	Info   					string `json:"1. Information"`
	Symbol					string `json:"2. Symbol"`
	Updated					string `json:"3. Last Refreshed"`
	OutputSize				string `json:"4. Output Size"`
	TimeZone				string `json:"5. Time Zone"`
}
type TimeSeriesEntry struct {
	Open   				string `json:"1. open"`
	High				string `json:"2. high"`
	Low					string `json:"3. low"`
	Close				string `json:"4. close"`
	Volume				string `json:"5. volume"`
}

type TimeSeries struct {
	MetaData 			*MetaData
	TimeSeriesEntries	[]TimeSeriesEntry
}

type jsonMap map[string]TimeSeriesEntry

type Response struct {
	MetaData 			*MetaData
	TimeSeries 			*TimeSeries
}
