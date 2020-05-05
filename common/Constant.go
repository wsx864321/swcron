package common

const(
	//返回成功
	RESP_OK = 200
	//通用错误码
	RESP_FAIL = 201
	//存储事件
	JOB_SAVE_EVENT = 1
	//删除事件
	JOB_DEL_EVEVT  = 2
	// 任务保存目录
	JOB_SAVE_DIR = "/cron/jobs/"
	//杀死任务目录
	JOB_LILL_DIR = "/cron/kill/"
)
