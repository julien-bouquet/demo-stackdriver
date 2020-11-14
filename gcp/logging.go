package gcp

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// Group logs in generic task by LoggerName
var envKeyLoggerName = "LOGGER_NAME"

// Set Job in Resource of logs
var envKeyLoggerJob = "LOGGER_JOB"

type loggerGenericTask struct {
	LoggerType string
	LoggerJob  string
	TaskID     string
	Client     *logging.Logger
}

func newLoggerGenericTask(client *logging.Client, taskID string) loggerGenericTask {
	logger := loggerGenericTask{}
	logger.LoggerType = "generic_task"
	loggerJob := os.Getenv(envKeyLoggerJob)
	logger.LoggerJob = loggerJob
	logger.TaskID = taskID
	loggerName := os.Getenv(envKeyLoggerName)
	logger.Client = client.Logger(loggerName)
	return logger
}

// CreateClientLogger aa
func CreateClientLogger(ctx context.Context) *logging.Client {
	projectID := getProjectID(ctx)
	client, _ := logging.NewClient(ctx, projectID)

	return client
}

// InitializeLogger aa
func InitializeLogger(ctx context.Context, client *logging.Client, taskID string) context.Context {
	loggerGenericTask := newLoggerGenericTask(client, taskID)
	return context.WithValue(ctx, "logger", loggerGenericTask)
}

func (logger loggerGenericTask) log(message string, context interface{}, severity logging.Severity) {
	payload := map[string]interface{}{}
	payload["message"] = message
	payload["context"] = context

	entry := logging.Entry{
		Payload: payload,
		Resource: &monitoredres.MonitoredResource{
			Type: logger.LoggerType,
			Labels: map[string]string{
				"job":     logger.LoggerJob,
				"task_id": logger.TaskID,
			},
		},
		Severity: logging.ParseSeverity("Debug"),
	}
	fmt.Println(payload)
	logger.Client.Log(entry)
}

func Debug(ctx context.Context, message string, context interface{}) {
	logger := ctx.Value("logger").(loggerGenericTask)
	logger.log(message, context, logging.ParseSeverity("Debug"))
}

func Warn(ctx context.Context, message string, context interface{}) {
	logger := ctx.Value("logger").(loggerGenericTask)
	logger.log(message, context, logging.ParseSeverity("Warning"))
}
