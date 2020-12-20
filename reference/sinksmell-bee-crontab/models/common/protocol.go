package common

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorhill/cronexpr"
)

// Job 任务结构
type Job struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
	TimeOut  string `json:"time_out"` // 神坑  go中解析前端传来的json  这里使用int的话 死活得不到值
}

// JobSchedulerPlan 任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 // 要调度的任务
	Expr     *cronexpr.Expression //解析好的 要执行的cronExpr
	NextTime time.Time            // 下次调度时间
}

// JobExecInfo 任务执行信息
type JobExecInfo struct {
	Job        *Job               // 正在执行的任务
	PlanTime   time.Time          //	计划调度时间
	RealTime   time.Time          //	实际调度时间
	CancelCtx  context.Context    //用于取消任务的上下文
	CancelFunc context.CancelFunc //取消方法
}

// JobExecResult 任务执行结果
// 开始调度的时间 与开始执行的时间是不一样的
// 调度时间之差 反映了调度器的效率
// 执行开始与结束时间之差 为程序运行时间
type JobExecResult struct {
	Type      int          // 结果类型 正常执行 kill 超时终止
	ExecInfo  *JobExecInfo // 执行状态
	Output    []byte       // 输出结果
	Err       error        // 错误信息
	StartTime time.Time    // 开始运行时间
	EndTime   time.Time    // 结束时间
}

// LogBuffer 日志缓存 批量插入任务日志 提交吞吐效率
// 当buffer满了或者定时器时间到了 执行插入操作
type LogBuffer struct {
	Logs []interface{} // 任务日志集合
}

// JobEvent 任务变化事件
type JobEvent struct {
	EventType uint //事件类型 在常量中有定义
	Job       *Job // 任务
}

// Response 通用的返回类型
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResponse 生产response 的方法
func NewResponse(code int, msg string, data interface{}) Response {
	return Response{code, msg, data}
}

// NewJobEvent 构造任务变化时间
func NewJobEvent(eType uint, job *Job) *JobEvent {
	return &JobEvent{eType, job}
}

// NewJobSchedulerPlan 构造执行计划
func NewJobSchedulerPlan(job *Job) (plan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		log.Errorf("parse cron expr err: %v", err)
		return
	}
	// 构造调度计划对象
	plan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

// NewJobExecInfo 构造一个执行状态
func NewJobExecInfo(plan *JobSchedulerPlan) (info *JobExecInfo) {
	// 创建可以取消的上下文
	ctx, cancelFunc := context.WithCancel(context.TODO())
	info = &JobExecInfo{
		Job:        plan.Job,
		PlanTime:   plan.NextTime,
		RealTime:   time.Now(),
		CancelCtx:  ctx,
		CancelFunc: cancelFunc,
	}
	return
}

// String jobEvent 的toString方法
func (event *JobEvent) String() string {
	return fmt.Sprintf("%d %+v\n", event.EventType, *(event.Job))
}

/*
	ExecInfo  *JobExecInfo // 执行状态
	Output    []byte       // 输出结果
	Err       error        // 错误信息
	StartTime time.Time    // 开始运行时间
	EndTime   time.Time    // 结束时间

*/
//JobExecResult的执行结果
func (result *JobExecResult) String() string {
	return fmt.Sprintf("job_name:%s\ttype:%s\noutput:%serr:%v",
		result.ExecInfo.Job.Name,
		CodeMessage(result.Type),
		string(result.Output),
		result.Err)
}

//
func (plan *JobSchedulerPlan) String() string {
	return fmt.Sprintf("id:%s, expr:%s, next_time:%s", plan.Job.ID, plan.Job.CronExpr, plan.NextTime)
}
