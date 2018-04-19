package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron" //why does it not work when i put an underscore infront of the statement like woith the mysql driver?
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/database"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"github.com/sevlyar/go-daemon"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/rpc"
	"fmt"
	"strconv"
	"time"
	"os"
)

// 1. execute prelminiary tasks, spawn childprocess as clone, kill parent process -> child runs as daemon
// 2. request new data from XYZ apis every X minutes and push them into the db. -> cronjobs
// 3. wait for http requests -> rest api kevins job
// 4. IPC via rpc -> child process waits listens for commands at port 12449
func main() {

	//Command Line Interface:
	//If no arguements are passed, the bot ist started
	//other available areguements: Status, Stop, Restart

	var osArgs []string
	for idx := range os.Args{
		osArgs = append(osArgs, os.Args[idx])
	}
	if len(osArgs) < 2  {
		osArgs = append(osArgs, "-help")
	}

	cmd := osArgs[1]

	switch cmd {
		case rpc.RPC_STOP: {
			c := rpc.MakeClient()
			c.Stop()
			os.Exit(0)
		}
		case "-help": {
			fmt.Println("Usage ./bot <command> \nCommands: -start, -stop, -help, - status, -update")
			os.Exit(0)
		}
		case rpc.RPC_STATUS: {
			c := rpc.MakeClient()
			c.Status()
			os.Exit(0)
		}
		case rpc.RPC_UPDATE_DB: {
			c := rpc.MakeClient()
			c.UpdateOnDemand()
			os.Exit(0)
		}
		case rpc.RPC_START: {
			//Pre
			database.CredentialsFromStoredCredentials()

			context := &daemon.Context{
				PidFileName: "pid",
				//PidFilePerm: 0644,
				//LogFileName: "log",
				//LogFilePerm: 0640,
				//WorkDir:     "./",
				//Umask:       027,
				Args: []string{utility.CRAWLER_PS_NAME, cmd},
			}

			child, _ := context.Reborn()

			if child != nil {
				//Post parent
				utility.LogToFile("Parent process start")
				fmt.Println("Daemon started; name: " + utility.CRAWLER_PS_NAME) //type ps -ef oder ps -aux to see daemon; deactivate with kill pid
			} else {
				defer context.Release()
				//Post child
				utility.LogToFile("Child process start; daemon name: " + utility.CRAWLER_PS_NAME + " args: " + os.Args[1])

				//Start rpc server an listen on port 12499
				rpcSvr := rpc.CrawlerServer{}
				go rpcSvr.StartRPCServer()

				//Make initial updates
				utility.LogToFile("Initial database updates")
				go database.UpdateCoinMarketCapCoinData()
				go database.UpdateCoinMarketCapGlobalMarketData()
				go database.UpdatePoloniexChartData()
				go database.UpdateBitInfoChartsBiggestWallets()
				go database.UpdateAplhaVantageTop20ETFs()
				go database.UpdateQuandlGoldChartData()
				utility.LogToFile("Initial database updates started")

				//Initialize a new cronjob for coinmarketcap API and Databse operations. Cron jobs are thread safe.
				utility.LogToFile("Initialize cronjobs")
				c := cron.New()
				c.AddFunc("@every 0h6m", database.UpdateCoinMarketCapCoinData) //less than 5:01 mins doesnt make sense since coinmarketcap only refreshes every 5 mins. will cause lots of duplicate db entry errors if frequency set below 5 mins
				c.AddFunc("@every 0h6m", database.UpdateCoinMarketCapGlobalMarketData)
				c.AddFunc("@every 4h5m", database.UpdatePoloniexChartData) // havent figured out yet how to decrease update period below 4h
				c.AddFunc("@every 0h30m", database.UpdateBitInfoChartsBiggestWallets)
				c.AddFunc("@daily", database.UpdateAplhaVantageTop20ETFs)
				c.AddFunc("@daily", database.UpdateQuandlGoldChartData)
				c.Start()
				defer c.Stop()
				utility.LogToFile(strconv.Itoa(len(c.Entries())) + " cronjobs initialized")
				//we need this loop because c.Start() return immideatly - just starts the cron job in its own routine.
				// Also this could later be used to regularily update a config read from a file. Currently status is dumped in to server log every 60 mins
				for {
					time.Sleep(60 * time.Minute)
					//not sure if the below 2 lines work as intended. further testing required.
					c := rpc.MakeClient()
					c.SilentStatus()
				}
			}
		}
		default: {
			fmt.Println("Illegal Parameter. Usage ./bot <command> \nCommands: -start, -stop, -help, -status, -update")
			os.Exit(0)
		}
	}
	utility.LogToFile("Parent process end")
}
