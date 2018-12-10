package main

import (
	"context"
	"net/http"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resp := c.Response()
			resp.Writer = NewCustomResponseWriter(resp.Writer)
			return next(c)
		}
	})
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	newGateway(ctx, e)
	glog.Infof("Starting listening at %s", ":8080")
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
	}
}
