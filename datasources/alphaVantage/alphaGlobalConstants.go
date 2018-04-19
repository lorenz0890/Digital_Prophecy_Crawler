package alphaVantage

const (
	API_KEY = 			"&apikey=XAFVR025CS661824"
	OUTPUTSIZE = 		"&outputsize=full"
	STOCKS_DAILY_FULL = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY" // &symbol=MSFT&outputsize=full&apikey=XAFVR025CS661824"
)

//ETF top 20 symbols based on http://etfdb.com/compare/volume/ on 05.03.2018 - orderd by max volume
// probably these should be read from cfg file for easier maintenance
const (
	STOCK_RANKING_SOURCE									= "http://etfdb.com/compare/volume/"
	SYMBOL_SPDR_SP500_ETF									= "&symbol=SPY"
	SYMBOL_ISHARES_MSCI_EM_ETF								= "&symbol=EMM"
	SYMBOL_FINANCIAL_SELECT_SECTOR_SPDR_FUND 				= "&symbol=XLF"
	SYMBOL_IPATH_SP500_VIX_SHORT_TERM_FUTURES_ETN			= "&symbol=VXX"
	SYMBOL_POWERSHARES_QQQ									= "&symbol=QQQ"
	SYMBOL_VANECK_VECTORS_GOLD_MINERS_ETF					= "&symbol=GDX"
	SYMBOL_PROSHARES_ULTRA_VIX_SHORT_TERM_FUTURES 			= "&symbol=UVXY"
	SYMBOL_VELOCITY_SHARES_DAILY2x_SHORT_TERM_ETN			= "&symbol=TVIX"
	SYMBOL_ISHARES_MSCI_EAFE_ETF							= "&symbol=EFA"
	SYMBOL_ISHARES_RUSSEL_2000_ETF							= "&symbol=IWM"
	SYMBOL_VELOCITY_SHARES_3X_LONG_CRUDE					= "&symbol=UWTI"
	SYMBOL_ISHARES_CHINA_LARGE_CAP_ETF						= "&symbol=FXI"
	SYMBOL_UTILITIES_SELECT_SECTOR_SPDR_FUND				= "&symbol=XLU"
	SYMBOL_ISHARES_MSCI_BRAZIL_ETF							= "&symbol=EWZ"
	SYMBOL_UNITED_STATES_OIL_FUND							= "&symbol=USO"
	SYMBOL_ALERIAN_MLP_ETF									= "&symbol=AMLP"
	SYMBOL_SPDR_SP_OIL_GAS_EXPLORATION_PRODUCTION_ETF		= "&symbol=XOP"
	SYMBOL_SPDR_BARCLAYS_HIGH_YIELD_BOND_ETF				= "&symbol=JNK"
	SYMBOL_ISHARES_IBOXX_USD_HIGH_YIELD_CORPORATE_BOND_ETF 	= "&symbol=HYG"
	SYMBOL_ENERGY_SELECT_SECTOR_SPDR_FUND					= "&symbol=XLE"
)