package main

import (
	"flag"
	"fmt"

	_ "github.com/ZhiqinYang/grpc-gateway-demo/runtime"
)

var config = struct {
	authProxy    bool
	authEndpoint string
	authPrefix   string
}{}

func init() {
	flag.BoolVar(&config.authProxy, "auth_proxy", false, "proxy auth service")
	flag.StringVar(&config.authEndpoint, "auth_endpoint", "127.0.0.1:11111", "auth service endpoint")
	flag.StringVar(&config.authPrefix, "auth_http_prefix", "", "auth http proxy prefix")
	fmt.Printf("config => \t %+v \n", config)
}

func registerModules() {
	//	RegisterW

}
