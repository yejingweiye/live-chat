package main

import (
	"flag"
	"fmt"
	"live-chat/gateway"
	"os"
	"runtime"
	"time"
)

var (
	confFile string
)

func initArgs() {
	//str, _ := os.Getwd()
	//fmt.Printf("str",str)
	flag.StringVar(&confFile, "config", "./gateway/client/gateway.json", "where gateway.json is.")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// 初始化参数
	initArgs()
	initEnv()

	//  初始化配置
	if err = gateway.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 统计
	if err = gateway.InitStats(); err != nil {
		goto ERR
	}

	// 初始化连接管理器
	if err = gateway.InitConnMgr(); err != nil {
		goto ERR
	}

	// 初始化websocket服务器
	if err = gateway.InitWSServer(); err != nil {
		goto ERR
	}
	// TODO

	for {
		time.Sleep(1 * time.Second)
	}

	os.Exit(0)

ERR:
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)

}
