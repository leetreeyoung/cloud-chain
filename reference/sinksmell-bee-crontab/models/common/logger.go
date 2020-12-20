package common

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Logger mongoDB 日志存储器
type Logger struct {
	Client        *mongo.Client
	LogCollection *mongo.Collection // 任务日志集合
	LogStream     chan *JobExecLog  // 任务执行日志流
}

// HTTPJobLog 任务执行日志 界面展示
// 解析为标准时间的日志结构体
type HTTPJobLog struct {
	JobName      string `json:"job_name" `      //任务名
	Command      string `json:"command" `       //执行命令
	Err          string `json:"err" `           //错误信息
	Output       string `json:"output" `        //任务输出
	PlanTime     string `json:"plan_time" `     // 计划开始时间
	ScheduleTime string `json:"schedule_time" ` // 实际调度时间
	StartTime    string `json:"start_time" `    // 开始运行时间
	EndTime      string `json:"end_time" `      // 结束运行时间
}

// JobExecLog 任务执行日志 MongoDB 存储
type JobExecLog struct {
	JobID        string `json:"job_id" bson:"job_id"`
	JobName      string ` json:"job_name" bson:"job_name"`          //任务名
	Command      string ` json:"command" bson:"command"`            //执行命令
	Err          string `json:"err" bson:"err"`                     //错误信息
	Output       string `json:"output" bson:"output"`               //任务输出
	PlanTime     int64  `json:"plan_time" bson:"plan_time"`         // 计划开始时间 时间戳
	ScheduleTime int64  `json:"schedule_time" bson:"schedule_time"` // 实际调度时间
	StartTime    int64  `json:"start_time" bson:"start_time"`       // 开始运行时间
	EndTime      int64  `json:"end_time" bson:"end_time"`           // 结束运行时间
}

// InitLogger 初始化Logger的单例
func NewLogger(ctx context.Context, mongoURL string) (logger *Logger, err error) {
	var (
		client     *mongo.Client
		collection *mongo.Collection
	)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	if err = client.Connect(ctx); err != nil {
		log.Errorf("connect err", err)
		return
	}
	collection = client.Database("cron").Collection("log")

	// 初始化
	logger = &Logger{
		Client:        client,
		LogCollection: collection,
		LogStream:     make(chan *JobExecLog, 1024),
	}

	return
}

// 日志存储
func (logger *Logger) WriteLoop() {

	var (
		execLog     *JobExecLog        //待写入的日志
		buffer      *LogBuffer         // 日志缓冲区
		maxSize     = 128              // 缓冲最大容量
		commitTimer *time.Timer        // 提交定时器
		timeOut     = 10 * time.Second //超时时间
	)

	for {
		// 使用buffer 和定时器机制
		// 实现定时批量提交
		// 提高吞吐率
		// 减少I/O次数
		if commitTimer == nil {
			commitTimer = time.NewTimer(timeOut)
		}
		if buffer == nil {
			// 初始化缓冲载体
			buffer = &LogBuffer{
				Logs: make([]interface{}, 0),
			}
		}

		select {
		case execLog = <-logger.LogStream:
			// 有日志传来
			buffer.Logs = append(buffer.Logs, execLog)
			if len(buffer.Logs) >= maxSize {
				// log.Info("execLog buffer满了!")
				logger.saveLogs(buffer)
				buffer = nil
				commitTimer.Reset(timeOut)
			}
		case <-commitTimer.C:
			// log.Info("log存储定时器到期！")
			// 定时器到期
			if buffer != nil {
				// 保存日志
				logger.saveLogs(buffer)
				buffer = nil
			}
			commitTimer.Reset(timeOut)
		}
	}

}

// 批量保存日志
func (logger *Logger) saveLogs(buffer *LogBuffer) {
	if buffer == nil || len(buffer.Logs) == 0 {
		log.Info("no log need to save")
		return
	}
	_, err := logger.LogCollection.InsertMany(context.TODO(), buffer.Logs)
	if err != nil {
		log.Errorf("batch save logs err", err)
		return
	}
	log.Info("batch save log success")
}

// 日志读取
