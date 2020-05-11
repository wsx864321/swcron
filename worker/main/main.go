package main

import (
	"../../worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var(
	confDir string
)

func initArgs(){
	flag.StringVar(&confDir, "config","./worker.json","worker.json")
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
	err = worker.InitConfig(confDir)
	if err != nil {
		goto ERR
	}
	//初始化任务执行器
	worker.InitExcutor()
	//启动任务调度器
	worker.InitScheduler()

	//初始化etcd管理器
	err = worker.InitJobManager()
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