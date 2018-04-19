package rpc

import (
	"github.com/valyala/gorpc"
	//"fmt"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"fmt"
)


//return pointer to a running rpc client
func MakeClient() *CrawlerClient {
	c := &gorpc.Client{
		// TCP/IP address of the server.
		Addr: RPC_IP_ADRESS + RPC_PORT,
	}
	c.Start()

	var client = CrawlerClient{c}

	return &client
}

func (c *CrawlerClient) Stop() {
	// All client methods issuing RPCs are thread-safe and goroutine-safe,
	// i.e. it is safe to call them from multiple concurrently running goroutines.
	c.exec.Call(RPC_STOP)
	utility.LogToFile("Stop request sent")
}

func (c *CrawlerClient) UpdateOnDemand() {
	// All client methods issuing RPCs are thread-safe and goroutine-safe,
	// i.e. it is safe to call them from multiple concurrently running goroutines.
	c.exec.Call(RPC_UPDATE_DB)
	utility.LogToFile("Update on demand request sent")
}

func (c *CrawlerClient) Status() {
	// All client methods issuing RPCs are thread-safe and goroutine-safe,
	// i.e. it is safe to call them from multiple concurrently running goroutines.
	//c.exec.Call(RPC_STOP)
	resp, err := c.exec.Call(RPC_STATUS)
	utility.LogToFile("Status request sent")
	if err != nil {
		utility.LogToFile(fmt.Sprintf("Error when sending request to server: %s", err))
	}
	if resp != nil {
		if resp.(string) != RPC_STATUS {
			utility.LogToFile(fmt.Sprintf("Response from the server: %+v", resp))
			fmt.Println(resp)
		}
	}
}

func (c *CrawlerClient) SilentStatus() {
	// All client methods issuing RPCs are thread-safe and goroutine-safe,
	// i.e. it is safe to call them from multiple concurrently running goroutines.
	//c.exec.Call(RPC_STOP)
	resp, err := c.exec.Call(RPC_STATUS)
	utility.LogToFile("Silent status request sent")
	if err != nil {
		utility.LogToFile(fmt.Sprintf("Error when sending request to server: %s", err))
	}
	if resp != nil {
		if resp.(string) != RPC_STATUS {
			utility.LogToFile(fmt.Sprintf("Response from the server: %+v", resp))
			//fmt.Println(resp)
		}
	}
}