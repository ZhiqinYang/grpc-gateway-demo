package main

import (
	"context"
	"flag"
	"net/http"

	lm "github.com/ZhiqinYang/grpc-gateway-demo/middleware"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// echo.MiddlewareFunc
// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func main() {
	var listen string
	flag.StringVar(&listen, "listen", "0.0.0.0:8080", "listen address")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e := echo.New()
	e.Use(lm.CustomResponseWriterWrapper())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	registerModules()
	newGateway(ctx, e)
	glog.Infof("Starting listening at %s", ":8080")
	if err := e.Start(listen); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
	}
}
