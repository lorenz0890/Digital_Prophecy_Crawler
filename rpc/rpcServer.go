package rpc

import (
	"github.com/valyala/gorpc"
	"fmt"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"os"
	"time"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/database"
)

func delayedOsExit()() {
	time.Sleep(5000 *time.Millisecond)
	utility.LogToFile("Shutting down crawler")
	os.Exit(3)
}
func (cs *CrawlerServer) StartRPCServer () {
	utility.LogToFile("RPC server start")
	if len(cs.ip) < 8 {
		cs.ip = RPC_IP_ADRESS
	}
	if len(cs.port) < 1 {
		cs.port = RPC_PORT
	}
	cs.svr = &gorpc.Server{
		// Accept clients on this TCP address.
		Addr: RPC_PORT,

		// Echo handler - do XYZ dependent on request and return back the message we received from the client
		Handler: func(clientAddr string, request interface{}) interface{} {
			utility.LogToFile(fmt.Sprintf("Obtained request %+v from the client %s\n", request, clientAddr))

			req := fmt.Sprintf("%+v", request)
			switch req  {
				case RPC_STOP : {
					go delayedOsExit()	//we have to do it this way because defer causes problems with the response
					return request
				}
				case RPC_STATUS : {
					var answear string
					answear += "Process name: " + utility.CRAWLER_PS_NAME +"\n"

					if database.TestDB() {
						answear += "DB is online at " + database.HOST_NAME + ".\n"
					} else { answear += "DB does not respons as expected. Check logfile for details.\n" }

					if database.TestDataSoures() {
						answear += "All datasources are responsive."
					} else {
						answear += "At least one datasource does not respond as expected. Check logfile for details.."
					}
					return answear
				}
				case RPC_UPDATE_DB : {
					database.UpdateAllOnDemand()
					return request
				}

			}
			return request
		},
	}
	if err := cs.svr.Serve(); err != nil {
		utility.LogToFile(fmt.Sprint("Cannot start rpc server: %s", err))
	}
}