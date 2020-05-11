package worker

import (
	"../common"
	"os/exec"
	"context"
	"time"
)
//执行器
type Excutor struct {

}

var (
	G_Excutor *Excutor
)

func (excute *Excutor)ExcuteJob(job *common.JobExcuteInfo) {
	var (
		cmd *exec.Cmd
		ret []byte
		err error
		excuteJobRet *common.ExcuteJobResult
	)
	//启动协程执行任务
	go func() {
		cmd = exec.CommandContext(context.TODO(),"/bin/bash","-c",job.Job.Command)
		//执行任务开始时间
		excuteJobRet.StartTime = time.Now()
		ret,err = cmd.CombinedOutput()
		//执行任务结束时间
		excuteJobRet.EndTime = time.Now()
		excuteJobRet.Output = ret
		excuteJobRet.Err = err
		//推送到任务返回结果管道
		G_Scheduler.push2ExcuteJobResult(excuteJobRet)
	}()
}

func InitExcutor(){
	G_Excutor = &Excutor{

	}
}