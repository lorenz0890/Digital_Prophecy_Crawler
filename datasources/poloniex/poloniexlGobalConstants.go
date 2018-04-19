package poloniex

const (
	CHART_DATA_URL    		= "https://poloniex.com/public?command=returnChartData" //&currencyPair=USDT_BTC&start=1405699200&end=9999999999&period=14400"
	CHART_DATA_START  		= "&start=1420070400"                                          // 01.01.2015 00:00
	CHART_DATA_END	  		= "&end=9999999999"
	CHART_DATA_PERIOD 		= "&period=14400"             	                              // 4 hours
	CHART_DATA_PERIOD_INT	= 14400
)

const ( // other currencies from top 15 coinmarket cap not listed at poloniex
	USDT_BTC  = "&currencyPair=USDT_BTC"
	USDT_XMR  = "&currencyPair=USDT_XMR"
	USDT_ETC  = "&currencyPair=USDT_ETC"
	USDT_LTC  = "&currencyPair=USDT_LTC"
	USDT_ETH  = "&currencyPair=USDT_ETH"
	USDT_XRP  = "&currencyPair=USDT_XRP"
	USDT_BCH  = "&currencyPair=USDT_BCH"
	USDT_DASH = "&currencyPair=USDT_DASH"
)