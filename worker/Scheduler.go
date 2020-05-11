package worker

import (
	"../common"
	"time"
)

type Scheduler struct {
	JobEvent        chan *common.JobEvent               //etcd中的任务调度队列
	JobPlanTable    map[string]*common.JobSchedulerPlan //任务调度计划表
	JobExcuteTable  map[string]*common.JobExcuteInfo    //正在执行的任务
	ExcuteJobResult chan *common.ExcuteJobResult         //执行任务后的返回结果
}

var (
	G_Scheduler *Scheduler
)

func InitScheduler(){
	G_Scheduler = &Scheduler{
		JobEvent:        make(chan *common.JobEvent, 1000),
		JobPlanTable:    make(map[string]*common.JobSchedulerPlan),
		JobExcuteTable:  make(map[string]*common.JobExcuteInfo),
		ExcuteJobResult: make(chan *common.ExcuteJobResult, 1000),
	}

	go G_Scheduler.scheduleLoop()
}
//尝试执行任务
func (scheduler *Scheduler)TryStartJob(plan *common.JobSchedulerPlan){
	var(
		jobExcuteInfo *common.JobExcuteInfo
		ok bool
	)
	//任务可能执行的很久，1分钟可能会调度60次，但是任务没有结束就只执行一次
	if _,ok = scheduler.JobExcuteTable[plan.Job.Name];ok{
		//任务还在执行队列中
		return
	}
	//构建执行任务
	jobExcuteInfo = common.BuildJobExcuteInfo(plan)
	//存入到正在执行的任务队列
	scheduler.JobExcuteTable[plan.Job.Name] = jobExcuteInfo
	//执行任务
	G_Excutor.ExcuteJob(jobExcuteInfo)
}

//计算下次调度时间
func (scheduler *Scheduler)TryScheduler() time.Duration {
	var (
		item *common.JobSchedulerPlan
		nextTime time.Time
		now time.Time
		diff time.Duration
	)
	//没有任务的时候先睡眠1s
	if len(scheduler.JobPlanTable) == 0 {
		return 1 * time.Second
	}
	now = time.Now()
	diff = 0
	for _,item = range scheduler.JobPlanTable {
		nextTime = item.NextTime
		//执行任务并更新下次
		if nextTime.Before(now) || nextTime.Equal(now){
			scheduler.TryStartJob(item)
			//fmt.Println("jobname",item.Job.Name," nextime:",item.NextTime)
			//更新下次执行时间
			item.NextTime = item.Expr.Next(nextTime)
		}
		//统计最近的要执行任务的时间，也就是睡眠时间
		if diff == 0 || diff > nextTime.Sub(now) {
			diff = nextTime.Sub(now)
		}
	}

	return diff
}

//协程调度任务
func (scheduler *Scheduler)scheduleLoop(){
	var (
		jobEvent  *common.JobEvent
		timer     *time.Timer
		sleepTime time.Duration
		excuteJobResult *common.ExcuteJobResult
	)

	sleepTime = scheduler.TryScheduler()
	timer = time.NewTimer(sleepTime)

	for {
		select {
		//监听任务的增删改
		case jobEvent = <- scheduler.JobEvent:
			scheduler.handleJobEvent(jobEvent)
		//最近的任务到期了，需要去执行任务，并且更新任务执行列表
		case <- timer.C:
		//监听任务执行结果
		case excuteJobResult = <- scheduler.ExcuteJobResult:
			//todo 收集任务执行结果
		}
		//等待下个任务的执行
		timer.Reset(scheduler.TryScheduler())
	}
}

//处理任务事件
func (scheduler *Scheduler)handleJobEvent(jobEvent *common.JobEvent){
	var(
		err error
		jobSchedulerPlan *common.JobSchedulerPlan
		ok bool
	)

	switch jobEvent.EventType{
	//增加cron
	case common.JOB_SAVE_EVENT:
		jobSchedulerPlan,err = common.BuildJobSchedulerPlan(jobEvent.Job)
		if err != nil {
			return
		}
		scheduler.JobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	//删除cron
	case common.JOB_DEL_EVEVT:
		if  jobSchedulerPlan,ok = scheduler.JobPlanTable[jobEvent.Job.Name];ok{
			delete(scheduler.JobPlanTable,jobEvent.Job.Name)
		}
	}
}

//push到管道JobEvent
func (scheduler *Scheduler)Push2JobEvent(jobEvent *common.JobEvent){
	scheduler.JobEvent <- jobEvent
}

//push任务执行的返回结果
func (scheduler *Scheduler)push2ExcuteJobResult(result *common.ExcuteJobResult){
	scheduler.ExcuteJobResult <- result
}

