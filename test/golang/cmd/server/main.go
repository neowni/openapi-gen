package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"test/generated/message"
	"test/generated/server"
)

func main() {
	engine := gin.Default()
	s := server.New(engine)

	//																			接口注册

	s.TestTag1.Op1(func(ctx context.Context, uri *message.TestTag1Op1URI, qry *message.TestTag1Op1Qry, req *message.TestTag1Op1Req) (rsp200 *message.TestTag1Op1Rsp200, err error) {
		rsp200 = new(message.TestTag1Op1Rsp200)

		rsp200.URI1 = uri.URI1
		rsp200.URI2 = uri.URI2

		rsp200.Qry1 = qry.Qry1
		rsp200.Qry2 = qry.Qry2
		rsp200.Qryo = qry.Qryo

		rsp200.Req1 = req.Req1
		rsp200.Req2 = req.Req2

		return rsp200, nil
	})

	s.TestTag1.Op2(func(ctx context.Context, uri *message.TestTag1Op2URI, qry *message.TestTag1Op2Qry, req *message.TestTag1Op2Req) (rsp200 *message.TestTag1Op2Rsp200, err error) {
		rsp200 = req

		return rsp200, nil
	})

	s.TestTag1.Op3(func(ctx context.Context, uri *message.TestTag1Op3URI, qry *message.TestTag1Op3Qry, req *message.TestTag1Op3Req) (rsp200 *message.TestTag1Op3Rsp200, err error) {
		rsp200 = req
		return rsp200, nil
	})

	s.TestTag2.Op4(func(ctx context.Context) (rsp200 *message.TestTag2Op4Rsp200, err error) {
		return new(message.TestTag2Op4Rsp200), nil
	})

	s.TestTag2.Op5(func(ctx context.Context) (rsp204 *message.TestTag2Op5Rsp204, err error) {
		return new(message.TestTag2Op5Rsp204), nil
	})

	s.TestTag2.Op6(func(ctx context.Context, req *message.TestTag2Op6Req) (rsp200 *message.TestTag2Op6Rsp200, rsp201 *message.TestTag2Op6Rsp201, rsp202 *message.TestTag2Op6Rsp202, rsp203 *message.TestTag2Op6Rsp203, err error) {
		switch req.Code {
		case 200:
			return new(message.TestTag2Op6Rsp200), nil, nil, nil, nil
		case 201:
			return nil, new(message.TestTag2Op6Rsp201), nil, nil, nil
		case 202:
			return nil, nil, new(message.TestTag2Op6Rsp202), nil, nil
		case 203:
			return nil, nil, nil, new(message.TestTag2Op6Rsp203), nil
		}

		return nil, nil, nil, nil, fmt.Errorf("op6 code %d", req.Code)
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:30435",
		Handler: engine,
	}

	srvChan := make(chan error, 16)
	go func() {
		srv.ListenAndServe()
		fmt.Println("服务器开始运行")
	}()

	// 等待信号
	quit := make(chan os.Signal, 16)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvChan:
		fmt.Println("服务停止运行", err)
	case <-quit:
		fmt.Println("接收到停止信号")
		err := srv.Shutdown(context.Background())
		fmt.Println("服务停止运行", err)
	}
}
