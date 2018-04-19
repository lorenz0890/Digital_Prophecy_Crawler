package quandl

type Response struct {
	SourceCode 	string       `json:"source_code"`
	SourceName 	string       `json:"source_name"`
	Code       	string       `json:"code"`
	Frequency  	string       `json:"frequency"`
	FromDate  	string       `json:"from_date"`
	ToDate    	string       `json:"to_date"`
	Columns   	[]string   	 `json:"column_names"`
	Data       	interface{}  `json:"data"`
}

type GoldDataEntry struct {
	Date		string       `json:"0"`
	UsdAm		float64      `json:"1"`
	UsdPm		float64      `json:"2"`
	GbpAm		float64      `json:"3"`
	GbpPm		float64      `json:"4"`
	EurAm		float64      `json:"5"`
	EurPm		float64      `json:"6"`
}

type GoldDataTimeSeries struct {
	GoldDataEntries []GoldDataEntry
}

type Quandl struct {
	GoldDataTimeSeries GoldDataTimeSeries
}