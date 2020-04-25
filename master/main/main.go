package main

import (
	"../../master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var(
	confDir string
)

func initArgs(){
	flag.StringVar(&confDir, "config","./master.json","指定master.json")
}

func initEnv(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
	var (
		err error
	)
	//获取初始化参数
	initArgs()
	//初始化线程数
	initEnv()
	//初始化配置
	err = master.InitConfig(confDir)
	if err != nil {
		goto ERR
	}
	//初始化etcd管理器
	err = master.InitJobManager()
	if err != nil {
		goto ERR
	}
	//启动httpserver
	err = master.InitServer()
	if err != nil {
		goto ERR
	}
	for {
		time.Sleep(time.Second)
	}

	return

	ERR:
		fmt.Println(err)
}