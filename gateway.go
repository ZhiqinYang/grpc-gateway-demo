package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
)

type module struct {
	prefix      string
	endpoint    string
	factory     factory
	middlewares []echo.MiddlewareFunc
	opts        []grpc.DialOption
}

type factory func(context.Context, *gwruntime.ServeMux, string, []grpc.DialOption) error

var (
	modules   []module
	marshaler = jsonpb.Marshaler{}
	pools     = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

// 注册网关
func RegisterGW(endpoint, prefix string, f factory, middlewares []echo.MiddlewareFunc, opts ...grpc.DialOption) {
	modules = append(modules, module{factory: f, endpoint: endpoint, opts: opts, prefix: prefix, middlewares: middlewares})
}

func newGateway(ctx context.Context, e *echo.Echo) {
	fn := gwruntime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, pb proto.Message) error {
		var buf = pools.Get().(*bytes.Buffer)
		buf.Reset()
		if err := marshaler.Marshal(buf, pb); err != nil {
			pools.Put(buf)
			return err
		}

		bts, err := json.Marshal(struct {
			Code int             `json:"code"`
			Msg  string          `json:"msg"`
			Data json.RawMessage `json:"data,omitempty"`
		}{
			Code: 200,
			Msg:  "Success",
			Data: buf.Bytes(),
		})

		// 回收数据 一定要放在marshal后面
		pools.Put(buf)
		if err != nil {
			return err
		}

		_, err = w.Write(bts)
		if err != nil {
			// 使用完后不允许再往body 填写数据
			((w.(*echo.Response).Writer).(*CustomResponseWriter)).Complete()
		}
		return err
	})

	for _, m := range modules {
		mux := gwruntime.NewServeMux(fn)
		checkErr(m.factory(ctx, mux, m.endpoint, m.opts))
		e.Any(fmt.Sprintf("%s/*", m.prefix), echo.WrapHandler(mux), m.middlewares...)
	}
}
