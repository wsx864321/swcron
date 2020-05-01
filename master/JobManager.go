package master
import (
	"../common"
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"time"
)
//etcd管理器结构体
type JobManager struct {
	CLient *clientv3.Client
	Kv     clientv3.KV
	Lease  clientv3.Lease
}

var (
	G_JobManager *JobManager
)
//初始化etcd管理器
func InitJobManager() error {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		lease  clientv3.Lease
		kv     clientv3.KV
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

	G_JobManager = &JobManager{
		CLient:client,
		Kv:kv,
		Lease:lease,
	}

	return nil
}
//保存job
func (jobManager *JobManager)SaveJob(job *common.Job) (*common.Job,error) {
	var (
		jobKey string
		jobStr []byte
		putRet *clientv3.PutResponse
		oldJob common.Job
		err    error
	)
	//job的完整key
	jobKey = common.JOB_SAVE_DIR + job.Name
	//json序列化
	jobStr,err = json.Marshal(job)
	if err != nil {
		return nil,err
	}
	//保存job任务
	putRet,err = jobManager.Kv.Put(context.TODO(), jobKey, string(jobStr),clientv3.WithPrevKV())
	//判断是否是对job进行修改
	if putRet.PrevKv != nil {
		err = json.Unmarshal(putRet.PrevKv.Value, &oldJob)
		if err != nil {
			return nil,err
		}
	}

	return &oldJob,err
}
//删除定时任务
func (jobManager *JobManager)DelJob(jobName string) (*common.Job,error) {
	var (
		err    error
		jobKey string
		delRet *clientv3.DeleteResponse
		oldJob common.Job
	)
	//job的完整key
	jobKey = common.JOB_SAVE_DIR + jobName
	//删除定时任务
	delRet,err = jobManager.Kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil,err
	}

	if len(delRet.PrevKvs) == 0 {
		return nil,errors.New("Jobs that do not exist")
	}

	err = json.Unmarshal(delRet.PrevKvs[0].Value, &oldJob)
	if err != nil {
		return nil,err
	}

	return &oldJob,nil
}
//获取所有job
func (jobManager *JobManager)JobList() ([]*common.Job, error) {
	var (
		getRet  *clientv3.GetResponse
		err     error
		job     *common.Job
		jobList []*common.Job
	)

	getRet,err = jobManager.Kv.Get(context.TODO(),common.JOB_SAVE_DIR,clientv3.WithPrefix())
	if err != nil {
		return nil,err
	}

	jobList = make([]*common.Job,0)
	for _,v := range getRet.Kvs {
		job = new(common.Job)
		err = json.Unmarshal([]byte(v.Value), job)
		if err != nil {
			continue
		}
		jobList = append(jobList, job)
	}

	return jobList,nil
}
//杀死任务
func (jobManager *JobManager)KillJob(jobName string) error {
	var(
		jobKey   string
		leaseRet *clientv3.LeaseGrantResponse
		err      error
		leaseId  clientv3.LeaseID
	)

	jobKey = common.JOB_LILL_DIR+jobName
	//让worker进程去监听put操作，创造一个租约，自动过期即可
	leaseRet,err = jobManager.Lease.Grant(context.TODO(),1)
	if err != nil {
		return err
	}
	//租约ID
	leaseId = leaseRet.ID
	//设置killer的标记
	_,err = jobManager.Kv.Put(context.TODO(),jobKey,"",clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}

	return nil

}
