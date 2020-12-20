package common

import (
	"encoding/json"
	"strings"
)

// ExtractJobID 从etcd key中提取出Job ID
func ExtractJobID(key string) string {
	return strings.TrimPrefix(key, JobSavePath)
}

// ExtractKillerID 从etcd killer 的key中提取Job ID
func ExtractKillerID(key string) string {
	return strings.TrimPrefix(key, JobKillerPath)
}

// ExtarctWorkerIP 从etcd /cron/worker/ip 中获取 worker 的ip
func ExtarctWorkerIP(key string) string {
	return strings.TrimPrefix(key, JobWorkerPath)
}

// UnpackJob 反序列化得到Job
func UnpackJob(value []byte) (ret *Job, err error) {
	ret = &Job{}
	if err = json.Unmarshal(value, &ret); err != nil {
		return
	}
	return
}
