package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"net/http"
	"strings"
	"time"
)

//job的结构体
type Job struct {
	Name     string `json:"name"`      // name
	Command  string `json:"command"`   //shell命令
	CronExpr string `json:"cron_expr"` //cron表达式
}
//响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
//任务事件
type JobEvent struct {
	EventType int
	Job       *Job
}
//job的执行计划
type JobSchedulerPlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}
func ApiResponse(w http.ResponseWriter,code int,msg string,data interface{}) error {
	var (
		res Response
		ret []byte
		err error
	)

	res.Code = code
	res.Msg  = msg
	res.Data = data

	ret,err = json.Marshal(res)
	if err != nil {
		return err
	}
	//设置响应头
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(ret)
	if err != nil {
		return err
	}

	return nil
}
//反序列化Job
func UnPackJob(data []byte) (*Job,error) {
	var (
		job Job
		err error
	)
	err = json.Unmarshal(data,&job)
	if err != nil {
		return nil,err
	}
	return &job,nil
}
//从etcd的key中提取任务名
func GetJobName(job string) string {
	return strings.TrimPrefix(job, JOB_SAVE_DIR)
}
//任务事件，1.存储时间 2.删除事件
func BuildJobEvent(eventType int,job *Job) *JobEvent {
	return &JobEvent{
		EventType:eventType,
		Job:job,
	}
}

func BuildJobSchedulerPlan(job *Job) (*JobSchedulerPlan,error) {
	var (
		err error
		expr *cronexpr.Expression
	)
	expr,err = cronexpr.Parse(job.Command)
	if err != nil {
		return nil,err
	}

	return &JobSchedulerPlan{
		Job:job,
		Expr:expr,
		NextTime:expr.Next(time.Now()),
	},nil
}