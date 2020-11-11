package main

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

func createClientLogger(ctx context.Context) *logging.Client {
	projectID := getProjectID(ctx)
	client, _ := logging.NewClient(ctx, projectID)

	return client
}

func initializeLogging(ctx context.Context, client *logging.Client, taskID string) (context.Context, *logging.Logger) {
	ctx = context.WithValue(ctx, "loggerType", "generic_task")
	loggerJob := os.Getenv(envKeyLoggerJob)
	ctx = context.WithValue(ctx, "loggerJob", loggerJob)
	ctx = context.WithValue(ctx, "loggerTaskID", taskID)

	loggerName := os.Getenv(envKeyLoggerName)

	return ctx, client.Logger(loggerName)
}

func debug(ctx context.Context, logger *logging.Logger, message string, context interface{}) {
	loggerType := fmt.Sprintf("%v", ctx.Value("loggerType"))
	loggerJob := fmt.Sprintf("%v", ctx.Value("loggerJob"))
	loggerTaskID := fmt.Sprintf("%v", ctx.Value("loggerTaskID"))

	payload := map[string]interface{}{}
	payload["message"] = message
	payload["context"] = context

	entry := logging.Entry{
		Payload: payload,
		Resource: &monitoredres.MonitoredResource{
			Type: loggerType,
			Labels: map[string]string{
				"job":     loggerJob,
				"task_id": loggerTaskID,
			},
		},
		Severity: logging.ParseSeverity("Debug"),
	}
	fmt.Println(payload)
	logger.Log(entry)
}
