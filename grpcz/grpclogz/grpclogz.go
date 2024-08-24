package grpclogz

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"google.golang.org/grpc/grpclog"
)

func NewGRPCLoggerV2(l *slog.Logger) grpclog.LoggerV2 {
	return &Logger{l: l}
}

var _ grpclog.LoggerV2 = (*Logger)(nil)

type Logger struct {
	l *slog.Logger
}

func (l *Logger) V(level int) bool                            { return l.l.Enabled(context.Background(), slog.Level(level)) }
func (l *Logger) Info(args ...interface{})                    { l.l.Info(fmt.Sprint(args...)) }
func (l *Logger) Infoln(args ...interface{})                  { l.l.Info(fmt.Sprint(args...)) }
func (l *Logger) Infof(format string, args ...interface{})    { l.l.Info(fmt.Sprintf(format, args...)) }
func (l *Logger) Warning(args ...interface{})                 { l.l.Warn(fmt.Sprint(args...)) }
func (l *Logger) Warningln(args ...interface{})               { l.l.Warn(fmt.Sprint(args...)) }
func (l *Logger) Warningf(format string, args ...interface{}) { l.l.Warn(fmt.Sprintf(format, args...)) }
func (l *Logger) Error(args ...interface{})                   { l.l.Error(fmt.Sprint(args...)) }
func (l *Logger) Errorln(args ...interface{})                 { l.l.Error(fmt.Sprint(args...)) }
func (l *Logger) Errorf(format string, args ...interface{})   { l.l.Error(fmt.Sprintf(format, args...)) }
func (l *Logger) Fatal(args ...interface{})                   { l.l.Error(fmt.Sprint(args...)); os.Exit(1) }
func (l *Logger) Fatalln(args ...interface{})                 { l.l.Error(fmt.Sprint(args...)); os.Exit(1) }
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.l.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
