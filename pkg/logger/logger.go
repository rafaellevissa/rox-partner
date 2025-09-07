package logger

import "log"

type Logger struct{}

func New() *Logger                                       { return &Logger{} }
func (l *Logger) Infof(fmtStr string, v ...interface{})  { log.Printf("INFO: "+fmtStr, v...) }
func (l *Logger) Errorf(fmtStr string, v ...interface{}) { log.Printf("ERROR: "+fmtStr, v...) }
func (l *Logger) Info(msg string)                        { log.Printf("INFO: %s", msg) }
