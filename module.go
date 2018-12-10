package main

import (
	"flag"
	"fmt"

	_ "git.kunlun/KUNLUN-Hyper/gateway/runtime"
)

var config = struct {
	authProxy    bool
	authEndpoint string
	authPrefix   string
}{}

// 设置要代理的服务
func init() {
	flag.BoolVar(&config.authProxy, "auth_proxy", false, "proxy auth service")
	flag.StringVar(&config.authEndpoint, "auth_endpoint", "127.0.0.1:11111", "auth service endpoint")
	flag.StringVar(&config.authPrefix, "auth_http_prefix", "", "auth http proxy prefix")
	// TODO other rpc
	flag.Parse()
	if config.authProxy {
		/*
			RegisterGW(config.authEndpoint,
				config.authPrefix,
				oauthpb.RegisterOauthHandlerFromEndpoint,
				nil,
				grpc.WithInsecure())
		*/
	}
	fmt.Printf("config => \t %+v \n", config)
}
