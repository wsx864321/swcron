package master

import (
	"../common"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)


type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer//单例Server
)

func handleJobSave(w http.ResponseWriter,r *http.Request){
	var (
		data string
		err error
		job common.Job
		oldJob *common.Job
	)

	err = r.ParseForm()
	if err != nil {
		goto ERR
	}

	data = r.PostForm.Get("job")
	//反序列化
	err= json.Unmarshal([]byte(data), &job)
	if err != nil {
		goto ERR
	}
	oldJob,err = G_JobManager.SaveJob(&job)
	if err != nil {
		goto ERR
	}

	err = common.ApiResponse(w, common.RESP_OK, "success", oldJob)
	if err != nil {
		goto ERR
	}
	return
ERR:
	common.ApiResponse(w, common.RESP_FAIL,err.Error(),nil)
}
//删除任务
func handleJobDel(w http.ResponseWriter,r *http.Request){
	var (
		err error
		JobName string
		oldJob *common.Job
	)
	err = r.ParseForm()
	if err != nil {
		goto ERR
	}

	JobName = r.PostForm.Get("job_name")
	oldJob,err = G_JobManager.DelJob(JobName)
	if err != nil {
		goto ERR
	}

	err = common.ApiResponse(w, common.RESP_OK, "success", oldJob)
	if err != nil {
		goto ERR
	}

	return
ERR:
	common.ApiResponse(w, common.RESP_FAIL,err.Error(),nil)
}
//获取所有任务
func handleJobList(w http.ResponseWriter,r *http.Request){
	var (
		jobList []*common.Job
		err     error
		data    map[string]*common.Job
	)

	jobList,err = G_JobManager.JobList()
	if err != nil {
		goto ERR
	}

	data = make(map[string]*common.Job)
	for _,v := range jobList {
		data[v.Name] = v
	}

	err = common.ApiResponse(w, common.RESP_OK, "", data)
	if err != nil {
		goto ERR
	}
	return
ERR:
	common.ApiResponse(w, common.RESP_FAIL,err.Error(),nil)
}

func InitServer() error {
	var (
		mux *http.ServeMux
		lister net.Listener
		err error
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save",handleJobSave)
	mux.HandleFunc("/job/del",handleJobDel)
	mux.HandleFunc("/job/list",handleJobList)
	//监听端口
	lister,err = net.Listen("tcp",":"+strconv.Itoa(G_config.Port))
	if err != nil {
		return err
	}
	httpServer = &http.Server{
		ReadTimeout:time.Duration(G_config.ReadTimeout) *  time.Millisecond,
		WriteTimeout:time.Duration(G_config.WriteTimeout) *  time.Millisecond,
		Handler:mux,//http.ServeMux实现了Server.Handler接口
	}

	G_apiServer = &ApiServer{
		httpServer:httpServer,
	}
	//启动httpserver
	go httpServer.Serve(lister)

	return nil
}