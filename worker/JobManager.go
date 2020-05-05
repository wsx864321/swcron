package worker

import (
	"../common"
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)
//etcd管理器结构体
type JobManager struct {
	CLient *clientv3.Client
	Kv     clientv3.KV
	Lease  clientv3.Lease
	Wc     clientv3.Watcher
}

var (
	G_JobManager *JobManager
)

func (jobManage *JobManager)watchJobs() {
	var (
		err           error
		getRet        *clientv3.GetResponse
		kvItem        *mvccpb.KeyValue
		job           *common.Job
		startRevision int64
		watchChanRet  clientv3.WatchChan
		wcResponse    clientv3.WatchResponse
		event         *clientv3.Event
		jobName       string
		jobEvent      *common.JobEvent
	)

	getRet, err = jobManage.Kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return
	}
	//遍历全部的任务
	for _, kvItem = range getRet.Kvs {
		job, err = common.UnPackJob(kvItem.Value)
		if err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_SAVE_EVENT, job)
			//push到调度器当中
			G_Scheduler.Push2JobEvent(jobEvent)
		}
	}
	//对/cron/jobs目录下的key进行监听
	go func() {
		startRevision = getRet.Header.Revision + 1
		watchChanRet = jobManage.Wc.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix(), clientv3.WithRev(startRevision))
		for wcResponse = range watchChanRet{
			for _,event = range wcResponse.Events {
				switch event.Type {
				case mvccpb.PUT:
					job,err = common.UnPackJob(event.Kv.Value)
					if err != nil {
						continue
					}

					jobEvent = common.BuildJobEvent(common.JOB_SAVE_EVENT, job)
				case mvccpb.DELETE:
					jobName = common.GetJobName(string(event.Kv.Key))
					job = &common.Job{
						Name:jobName,
					}
					jobEvent = common.BuildJobEvent(common.JOB_DEL_EVEVT, job)
				}
				//push到调度器当中
				G_Scheduler.Push2JobEvent(jobEvent)
			}
		}
	}()
}


//初始化etcd管理器
func InitJobManager() error {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		lease  clientv3.Lease
		kv     clientv3.KV
		wc     clientv3.Watcher
	)
	config = clientv3.Config{
		Endpoints:[]string{G_config.EtcdEndpoints},
		DialTimeout:time.Duration(G_config.EtcdDialTimeout)*time.Millisecond,
	}

	client,err = clientv3.New(config)
	if err != nil {
		return err
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	wc = clientv3.NewWatcher(client)

	G_JobManager = &JobManager{
		CLient:client,
		Kv:kv,
		Lease:lease,
		Wc:wc,
	}

	//启动任务监听
	G_JobManager.watchJobs()

	return nil
}
