package common

const (
	// 任务事件
	// 任务保存事件
	JobEventSave = iota
	// 任务删除事件
	JobEventDelete
	// 杀死任务事件
	JobEventKill

	// 任务保存目录
	JobSavePath = "/cron/jobs/"
	// job killer 目录
	JobKillerPath = "/cron/killer/"
	// 分布式锁路径
	JobLockPath = "/cron/lock/"
	// worker节点注册路径  服务注册
	JobWorkerPath = "/cron/worker/"

	// 任务执行结果
	ResSuccess = 0 // 任务正常执行结束
	ResKilled  = 1 // 任务被提前终止
	ResTimeout = 2 // 任务超时自动终止

)

func CodeMessage(code int) string {

	var msg string
	switch code {
	case ResSuccess:
		msg = "success"
	case ResKilled:
		msg = "killed"
	case ResTimeout:
		msg = "timeout"
	}

	return msg
}
