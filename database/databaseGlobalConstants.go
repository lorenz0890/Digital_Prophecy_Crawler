package database

const (
	HOST_NAME        				= "crypto.cveqn4aeag4g.eu-central-1.rds.amazonaws.com:3306" // "159.89.8.104:3306" "127.0.0.1:3306"

	COINMARKETCAP_DB 				= "cryptoMaster"
	POLONIEX_DB 	 				= "cryptoMaster"// previously "poloniex_data"
	BITINFOCHARTS_DB 				= "cryptoMaster"// previously "bitinfocharts_data"
	ALPHAVANTAGE_DB	 				= "cryptoMaster"// previously "alphavantage_data"
	QUANDL_DB	 					= "cryptoMaster"

	COINMARKETCAP_COIN_TABLE 		= "coinmarketcapCoinData"
	COINMARKETCAP_GLOBAL_TABLE		= "coinmarketcapGlobalMarketData"
	POLONIEX_CHART_TABLE			= "poloniexChartData"
	BITINFOCHARTS_BIGWALLET_TABLE	= "bitinfochartsBigWallets"
	ALPHAVANTAGE_TOP20ETFs_TABLE	= "alphavantageTop20ETFs"
	QUANDL_GOLD_LBMA_FIXED_TABLE	= "quandlGoldChartDataLBMAFixed"

	DRIVER_NAME     				= "mysql"
	BOT_LOGIN_NAME  	 			= "bot"
	BOT_PASSWORD    				= "R4F6hYbn731Kf1Fw"
)
