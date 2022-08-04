package cohesivenet

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

/*
This is a logger that implements the CanLog interface defined in the cohesivenet
go client such that we can ensure the client properly provides terraform logs
*/

type Logger struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) Logger {
	return Logger{
		ctx: ctx,
	}
}

func (logger Logger) SetLevel(logLevel int) {}

func (logger Logger) DisableLogs() {}

func (logger Logger) Debug(message string) {
	tflog.Debug(logger.ctx, message)
}

func (logger Logger) Info(message string) {
	tflog.Info(logger.ctx, message)
}

func (logger Logger) Warn(message string) {
	tflog.Warn(logger.ctx, message)
}

func (logger Logger) Error(message string) {
	tflog.Error(logger.ctx, message)
}
