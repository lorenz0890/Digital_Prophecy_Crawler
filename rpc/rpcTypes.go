package rpc

import "github.com/valyala/gorpc"

type CrawlerClient struct {
	exec *gorpc.Client
}

type CrawlerServer struct {
	svr *gorpc.Server
	ip string
	port string
}

