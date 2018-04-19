package database

import (
	_ "github.com/go-sql-driver/mysql"
)

var dbCredentials DatabaseCredentials

type DatabaseCredentials struct {
	driverName              string
	dataSourceCoinMarketCap string
	dataSourcePoloniex      string
	dataSourceBitInfoCharts string
	dataSourceAlphaVantage	string
	dataSourceQuandl		string
}