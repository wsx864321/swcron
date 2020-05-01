package worker

import (
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
	"../common"
	"context"
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

func (jobManage *JobManager)watchJobs() error {
	var (
		err error
		getRet *clientv3.GetResponse
		kvItem *mvccpb.KeyValue
	)

	getRet,err = jobManage.Kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	//遍历全部的任务
	for _,kvItem = range getRet.Kvs {

	}
	
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

	return nil
}
