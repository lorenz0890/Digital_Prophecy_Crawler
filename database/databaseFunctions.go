package database

import (
	"database/sql"
	"strconv"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/datasources/coinMarketCap"
	"bufio"
	"os"
	"fmt"
	"reflect"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/datasources/poloniex"
	"strings"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/datasources/bitInfoCharts"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/datasources/alphaVantage"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/datasources/quandl"
)

//update on demand
func UpdateAllOnDemand(){
	go UpdateCoinMarketCapCoinData()
	go UpdateCoinMarketCapGlobalMarketData()
	go UpdatePoloniexChartData()
	go UpdateBitInfoChartsBiggestWallets()
	go UpdateAplhaVantageTop20ETFs()
	go UpdateQuandlGoldChartData()
	utility.LogToFile("Updating all tables on demand")
}

//checks if database is online
func TestDB () bool {
	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceCoinMarketCap)
	utility.CheckErr(err)
	if err == nil && db != nil { defer db.Close(); return true }
	if err == nil && db == nil { return false }
	if err != nil && db == nil { return false }
	return false
}

// check if datasources are OK. If one datasource fails, function returns false
func TestDataSoures () bool {
	cm := coinMarketCap.CoinMarketCap{}
	cm.GetCoinData(coinMarketCap.COIN_DATA__URL)
	if cm.CoinData == nil { return false }

	cm.GetGlobalMarketData(coinMarketCap.GLOBAL_MARKET_DATA__URL)
	if cm.GlobalMarketData.ActiveMarkets == 0 { return false }

	px := poloniex.Poloniex{}
	px.GetChartData(poloniex.CHART_DATA_URL + poloniex.USDT_BTC + poloniex.CHART_DATA_START + poloniex.CHART_DATA_END + poloniex.CHART_DATA_PERIOD)
	if px.ChartData == nil { return false }

	bw := bitInfoCharts.GetBiggestWallets()
	if bw == nil { return false }
	if len(bw) < 1 { return false }

	top := alphaVantage.GetTop20ETFs(alphaVantage.STOCKS_DAILY_FULL+alphaVantage.SYMBOL_SPDR_SP500_ETF+alphaVantage.OUTPUTSIZE+alphaVantage.API_KEY, alphaVantage.SYMBOL_SPDR_SP500_ETF)
	if top == nil { return false }
	if len(top) < 1 { return false }

	ql := quandl.Quandl{}
	ql.GetLMBAGold(quandl.QUANDL_GOLD_TIMESERIES + quandl.QUANDL_API_KEY + quandl.QUANDL_START_DATE)
	if ql.GoldDataTimeSeries.GoldDataEntries == nil { return false }
	if len(ql.GoldDataTimeSeries.GoldDataEntries) < 1 { return false}

	return true
}

func UpdateCoinMarketCapCoinData() {
	// open connection to database; defer its closure

	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceCoinMarketCap) // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		// get data from api
		cm := coinMarketCap.CoinMarketCap{}
		cm.GetCoinData(coinMarketCap.COIN_DATA__URL)

		for _, element := range cm.CoinData {
			s := "INSERT "+COINMARKETCAP_COIN_TABLE +" SET Id=?, Name=?, Symbol=?, Rank=?, PriceUsd=?, PriceBtc=?, Usd24hVolume=?, MarketCapUsd=?, AvailableSupply=?, TotalSupply=?, MaxSupply=?, PercentageChange1h=?, PercentageChange24h=?, PercentageChange7d=?, LastUpdated=?"
			stmt, err := db.Prepare(s)
			utility.CheckErr(err)
			if err == nil {
				//defer stmt.Close() //defer in loop dangerous -> oinly called when loop is completely finished

				i64, err := strconv.ParseInt(element.LastUpdated, 10, 64)
				i64 /= 10
				updated := time.Unix(i64*10^9, i64)

				maxSupply := element.MaxSupply
				if maxSupply == "" {
					maxSupply = "0"
				}

				res, err := stmt.Exec(element.Id, element.Name, element.Symbol, element.Rank, element.PriceUsd, element.PriceBtc, element.Usd24hVolume, element.MarketCapUsd, element.AvailableSupply, element.TotalSupply, maxSupply, element.PercentChange1h, element.PercentChange24h, element.PercentChange7d, updated)
				utility.CheckErr(err)
				if res != nil {
					_, err := res.LastInsertId()
					utility.CheckErr(err)
				}
				stmt.Close()
			}
		}
		utility.LogToFile("Updated Coinmarketcap global coin data")
	}
}

func UpdateCoinMarketCapGlobalMarketData() {
	// open connection to database; defer its closure
	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceCoinMarketCap) // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		// get data from api
		cm := coinMarketCap.CoinMarketCap{}
		cm.GetGlobalMarketData(coinMarketCap.GLOBAL_MARKET_DATA__URL)

		s := "INSERT "+ COINMARKETCAP_GLOBAL_TABLE +" SET TotalMarketCapUsd=?, Total24hVolumeUsd=?, BtcPercentOfMarket=?, ActiveCurrencies=?, ActiveAssets=?, ActiveMarkets=?, LastUpdated=?"
		stmt, err := db.Prepare(s)
		utility.CheckErr(err)
		if err == nil {
			defer stmt.Close()

			i64 := cm.GlobalMarketData.LastUpdated
			i64 /= 10
			updated := time.Unix(i64*10^9, i64)

			res, err := stmt.Exec(cm.GlobalMarketData.TotalMarketCapUsd, cm.GlobalMarketData.Total24hVolumeUsd, cm.GlobalMarketData.BitcoinPercentageOfMarketCap, cm.GlobalMarketData.ActiveCurrencies, cm.GlobalMarketData.ActiveAssets, cm.GlobalMarketData.ActiveMarkets, updated)
			utility.CheckErr(err)
			if res != nil {
				_, err := res.LastInsertId()
				utility.CheckErr(err)
			}
		}
		utility.LogToFile("Updated Coinmarketcap global market data")
	}
}

func UpdatePoloniexChartData() {
	// open connection to database; defer its closure
	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourcePoloniex + "?parseTime=true") // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		// get data from api
		currencyPairs := [8]string{
			poloniex.USDT_BTC,
			poloniex.USDT_BCH,
			poloniex.USDT_DASH,
			poloniex.USDT_ETC,
			poloniex.USDT_ETH,
			poloniex.USDT_LTC,
			poloniex.USDT_XMR,
			poloniex.USDT_XRP,
		}

		//attention - every trading pair nmeeds own starting time
		for _, pair := range currencyPairs {

			//check date of last entry
			pairPure := strings.Split(pair, "=")[1]

			s := "SELECT Date FROM "+POLONIEX_CHART_TABLE+" WHERE CurrencyPair = " + "'" + pairPure + "'" + " ORDER BY Date "
			res, err := db.Query(s)
			if err == nil {
				utility.CheckErr(err)
				//defer res.Close() defer in einer for schleife ist gefährlich -> resource leak

				var d time.Time
				for res.Next() {
					res.Scan(&d)
				}
				chartDataStart := "&start=" + strconv.FormatInt(d.Unix(), 10)
				res.Close()

				//get the data from poloniex
				px := poloniex.Poloniex{}
				if d.Unix() >= 0 {
					px.GetChartData(poloniex.CHART_DATA_URL + pair + chartDataStart + poloniex.CHART_DATA_END + poloniex.CHART_DATA_PERIOD)
					//utility.LogToFile(poloniex.CHART_DATA_URL + pair + chartDataStart + poloniex.CHART_DATA_END + poloniex.CHART_DATA_PERIOD)
				} else {
					px.GetChartData(poloniex.CHART_DATA_URL + pair + poloniex.CHART_DATA_START + poloniex.CHART_DATA_END + poloniex.CHART_DATA_PERIOD)
					//utility.LogToFile(poloniex.CHART_DATA_URL + pair + chartDataStart + poloniex.CHART_DATA_END + poloniex.CHART_DATA_PERIOD)
				}

				//push the data into our db
				for _, element := range px.ChartData {
					s := "INSERT "+POLONIEX_CHART_TABLE+" SET CurrencyPair=?, Date=?, High=?, Low=?, Open=?, Close=?, Volume=?, QuoteVolume=?"
					stmt, err := db.Prepare(s)
					utility.CheckErr(err)
					if err == nil {
						//defer stmt.Close() //possible resource leak if within for loop -> defer only called at end of loop, when we leave the loop

						i64 := element.Date
						i64 /= 10
						date := time.Unix(i64*10^9, i64)

						res, err := stmt.Exec(pairPure, date, element.High, element.Low, element.Open, element.Close, element.Volume, element.QuoteVolume)

						utility.CheckErr(err)
						if res != nil {
							_, err := res.LastInsertId()
							utility.CheckErr(err)
						}
						stmt.Close()
					}
				}
			}
		}
	}
	utility.LogToFile("Updated Poloniex chart data")
}

func UpdateBitInfoChartsBiggestWallets() {

	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceBitInfoCharts+ "?parseTime=true") // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		bw := bitInfoCharts.GetBiggestWallets()

		for _, elem := range bw {
			s := "INSERT "+BITINFOCHARTS_BIGWALLET_TABLE+" SET Adresse=?, AmountBTC=?, ChangeBTC=?, Updated=?"
			stmt, err := db.Prepare(s)
			if err == nil {
				utility.CheckErr(err)

				lastUpdated := strings.Split(time.Now().String(), ".")[0]

				res, err := stmt.Exec(elem.Address, elem.AmountBTC, elem.ChangeBTC, lastUpdated)
				utility.CheckErr(err)
				if err == nil {
					if res != nil {
						_, err := res.LastInsertId()
						utility.CheckErr(err)
					}
				}
				stmt.Close()
			}
		}
	}
	utility.LogToFile("Updated Bitfinforcharts biggest wallets")
}

func UpdateAplhaVantageTop20ETFs() {
	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceAlphaVantage + "?parseTime=true") // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		symbols := []string{
			alphaVantage.SYMBOL_SPDR_SP500_ETF,
			alphaVantage.SYMBOL_ALERIAN_MLP_ETF,
			alphaVantage.SYMBOL_ENERGY_SELECT_SECTOR_SPDR_FUND,
			alphaVantage.SYMBOL_FINANCIAL_SELECT_SECTOR_SPDR_FUND,
			alphaVantage.SYMBOL_IPATH_SP500_VIX_SHORT_TERM_FUTURES_ETN,
			alphaVantage.SYMBOL_ISHARES_CHINA_LARGE_CAP_ETF,
			alphaVantage.SYMBOL_ISHARES_IBOXX_USD_HIGH_YIELD_CORPORATE_BOND_ETF,
			alphaVantage.SYMBOL_ISHARES_MSCI_BRAZIL_ETF,
			alphaVantage.SYMBOL_ISHARES_MSCI_EAFE_ETF,
			alphaVantage.SYMBOL_ISHARES_MSCI_EM_ETF,
			alphaVantage.SYMBOL_ISHARES_RUSSEL_2000_ETF,
			alphaVantage.SYMBOL_POWERSHARES_QQQ,
			alphaVantage.SYMBOL_PROSHARES_ULTRA_VIX_SHORT_TERM_FUTURES,
			alphaVantage.SYMBOL_SPDR_BARCLAYS_HIGH_YIELD_BOND_ETF,
			alphaVantage.SYMBOL_SPDR_SP_OIL_GAS_EXPLORATION_PRODUCTION_ETF,
			alphaVantage.SYMBOL_UNITED_STATES_OIL_FUND,
			alphaVantage.SYMBOL_UTILITIES_SELECT_SECTOR_SPDR_FUND,
			alphaVantage.SYMBOL_VANECK_VECTORS_GOLD_MINERS_ETF,
			alphaVantage.SYMBOL_VELOCITY_SHARES_3X_LONG_CRUDE,
			alphaVantage.SYMBOL_VELOCITY_SHARES_DAILY2x_SHORT_TERM_ETN,
		}

		//find date of last entry
		var d time.Time
		for idxSymbols := range symbols {
			s := "SELECT Date FROM "+ALPHAVANTAGE_TOP20ETFs_TABLE+" WHERE Symbol = " + "'" + strings.Split(symbols[idxSymbols], "=")[1] + "'" + " ORDER BY Date "
			res, err := db.Query(s)
			utility.CheckErr(err)
			if err == nil {
				for res.Next() {
					res.Scan(&d)
				}
				res.Close()
			}
		}
		//fmt.Println(strconv.FormatInt(d.Unix(), 10))

		// get data from alphavantage and if date is bigger than date of last entry, push in db
		for idxSymbols := range symbols {
			top := alphaVantage.GetTop20ETFs(alphaVantage.STOCKS_DAILY_FULL+symbols[idxSymbols]+alphaVantage.OUTPUTSIZE+alphaVantage.API_KEY, symbols[idxSymbols])
			if top != nil {
				for idx1 := range top {
					if len(top[idx1]) > 0 {
						t, err := time.Parse("2006-01-02", top[idx1][0])
						utility.CheckErr(err)
						//only push thos from last month if date of last entry is positive (if db is not empty)
						if d.Unix() > 0 {
							if t.Unix() > time.Now().Unix()-604800 {
								s := "INSERT "+ALPHAVANTAGE_TOP20ETFs_TABLE+" SET Date=?, Symbol=?, Open=?, Close=?, High=?, Low=?, Volume=?"
								stmt, err := db.Prepare(s)
								utility.CheckErr(err)

								op, _ := strconv.ParseFloat(top[idx1][2], 64)
								cl, _ := strconv.ParseFloat(top[idx1][3], 64)
								hi, _ := strconv.ParseFloat(top[idx1][4], 64)
								lo, _ := strconv.ParseFloat(top[idx1][5], 64)
								vo, _ := strconv.ParseFloat(top[idx1][6], 64)

								res, err := stmt.Exec(top[idx1][0], top[idx1][1], op, cl, hi, lo, vo)
								utility.CheckErr(err)
								if res != nil {
									_, err := res.LastInsertId()
									utility.CheckErr(err)
								}
								stmt.Close()
							}
						} else {
							s := "INSERT "+ALPHAVANTAGE_TOP20ETFs_TABLE+" SET Date=?, Symbol=?, Open=?, Close=?, High=?, Low=?, Volume=?"
							stmt, err := db.Prepare(s)
							utility.CheckErr(err)

							op, _ := strconv.ParseFloat(top[idx1][2], 64)
							cl, _ := strconv.ParseFloat(top[idx1][3], 64)
							hi, _ := strconv.ParseFloat(top[idx1][4], 64)
							lo, _ := strconv.ParseFloat(top[idx1][5], 64)
							vo, _ := strconv.ParseFloat(top[idx1][6], 64)

							res, err := stmt.Exec(top[idx1][0], top[idx1][1], op, cl, hi, lo, vo)
							utility.CheckErr(err)
							if res != nil {
								_, err := res.LastInsertId()
								utility.CheckErr(err)
							}
							stmt.Close()
						}
					}
				}
			} else {
				utility.LogToFile("alphaVantage.GetTop20s(url, etf string) [][]string returned nil for etf = " + symbols[idxSymbols])
			}
		}
	}
	utility.LogToFile("Updated Alpha Vantage top 20 etfs")
}

func UpdateQuandlGoldChartData () {
	db, err := sql.Open(dbCredentials.driverName, dbCredentials.dataSourceQuandl) // we need to code a timout function using a different thread - currently server will hang if cant connect to db
	utility.CheckErr(err)
	if err == nil {
		defer db.Close()

		//find latest db entry
		s := "SELECT EntryDate FROM "+QUANDL_GOLD_LBMA_FIXED_TABLE +" ORDER BY EntryDate ASC "
		res, err := db.Query(s)
		utility.CheckErr(err)

		var dateLatestString string
		cnt := 0
		if err == nil {
			for res.Next() {
				cnt++
				res.Scan(&dateLatestString)
			}
			res.Close()
		}
		dateLatest, err := time.Parse("2006-01-02 00:00:00", dateLatestString)
		utility.CheckErr(err)

		//fmt.Println(dateLatest.String())

		//calc new startdate -  if table is empty then dateLatest date is 0001.01.1 00:00:00 which is a negative unix time, so we need to check dates if unix time is > 0
		if dateLatest.Unix() > 0 {
			// get data from api
			ql := quandl.Quandl{}
			ql.GetLMBAGold(quandl.QUANDL_GOLD_TIMESERIES + quandl.QUANDL_API_KEY + quandl.QUANDL_START_DATE)

			//push data into db
			for _, element := range ql.GoldDataTimeSeries.GoldDataEntries {
				s := "INSERT "+ QUANDL_GOLD_LBMA_FIXED_TABLE +" SET EntryDate=?, USD_AM=?, USD_PM=?, GBP_AM=?, GBP_PM=?, EUR_AM=?, EUR_PM=?"
				stmt, err := db.Prepare(s)
				utility.CheckErr(err)

				if err == nil {

					entryDate, err := time.Parse("2006-01-02", element.Date)
					utility.CheckErr(err)
					if err == nil{
						if entryDate.Unix() > dateLatest.Unix() {
							res, err := stmt.Exec(entryDate, element.UsdAm, element.UsdPm, element.GbpAm, element.GbpPm, element.EurAm, element.EurPm)
							utility.CheckErr(err)
							if res != nil {
								_, err := res.LastInsertId()
								utility.CheckErr(err)
							}
						}
					}
					stmt.Close()
				}
			}
		} else {
			dateStartString := strings.Split(quandl.QUANDL_START_DATE, "=")[1]
			dateStart, err := time.Parse("2006-01-02", dateStartString)
			utility.CheckErr(err)

			// get data from api
			ql := quandl.Quandl{}
			ql.GetLMBAGold(quandl.QUANDL_GOLD_TIMESERIES + quandl.QUANDL_API_KEY + quandl.QUANDL_START_DATE)

			//push data into db
			for _, element := range ql.GoldDataTimeSeries.GoldDataEntries {
				s := "INSERT "+ QUANDL_GOLD_LBMA_FIXED_TABLE +" SET EntryDate=?, USD_AM=?, USD_PM=?, GBP_AM=?, GBP_PM=?, EUR_AM=?, EUR_PM=?"
				stmt, err := db.Prepare(s)
				utility.CheckErr(err)

				if err == nil {

					entryDate, err := time.Parse("2006-01-02", element.Date)
					utility.CheckErr(err)

					if err == nil{
						if entryDate.Unix() >dateStart.Unix() {
							res, err := stmt.Exec(entryDate, element.UsdAm, element.UsdPm, element.GbpAm, element.GbpPm, element.EurAm, element.EurPm)
							utility.CheckErr(err)
							if res != nil {
								_, err := res.LastInsertId()
								utility.CheckErr(err)
							}
						}
					}
					stmt.Close()
				}
			}
		}
	}
	utility.LogToFile("Updated Quandl LBMA fixed gold data")
}

func CredentialsFromConsole() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter DB Username: ")
	username, _ := reader.ReadString('\n')
	sz := len(username)
	username = username[:sz-1]

	fmt.Print("Enter DB Password: ")
	password, _ := reader.ReadString('\n')
	sz = len(password)
	password = password[:sz-1]

	// this should be replaced by a make function to improve code readability
	dbCredentials = DatabaseCredentials{
		DRIVER_NAME,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + COINMARKETCAP_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + POLONIEX_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + BITINFOCHARTS_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + ALPHAVANTAGE_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + QUANDL_DB,
	}
}

//gets username and password from command line parameters
//login where? there will be multiple databases in the future, need to handel that as well. maybe a thrird command line arguement? or cfg file for storing login information.
func CredentialsFromParameters() {
	username := os.Args[1]
	password := os.Args[2]
	fmt.Println(username, len(username))
	fmt.Println(password, len(password))
	if reflect.TypeOf(username) == nil || reflect.TypeOf(password) == nil {
		fmt.Println("Usage: bot username password")
		os.Exit(-1)
	}

	// this should be replaced by a make function to improve code readability
	dbCredentials = DatabaseCredentials{
		DRIVER_NAME,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + COINMARKETCAP_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + POLONIEX_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + BITINFOCHARTS_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + ALPHAVANTAGE_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + QUANDL_DB,
	}
}

//login using stored credentials. fancy cmd line interface would use all 3 functions, depending on startup type
func CredentialsFromStoredCredentials() {

	// this should be replaced by a make function to improve code readability
	dbCredentials = DatabaseCredentials{
		DRIVER_NAME,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + COINMARKETCAP_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + POLONIEX_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + BITINFOCHARTS_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + ALPHAVANTAGE_DB,
		BOT_LOGIN_NAME + ":" + BOT_PASSWORD + "@tcp(" + HOST_NAME + ")/" + QUANDL_DB,
		}
}