package common

// JobFilter  根据任务名在MongoDB中查询日志的过滤器
type JobFilter struct {
	ID string `bson:"job_id"`
}
