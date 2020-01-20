package util

import "fmt"

type Logger struct {
}

func (l *Logger) Info(format string, val ...interface{}) {
	fmt.Printf("[info]"+format+"\n", val...)
}

func (l *Logger) Warn(format string, val ...interface{}) {
	fmt.Printf("[warn]"+format+"\n", val...)
}

func (l *Logger) Error(err error, format string, val ...interface{}) {
	msg := fmt.Sprintf("[error]"+format+"\n", val...)
	fmt.Printf(msg + err.Error() + "\n")
}

func (l *Logger) Debug(format string, val ...interface{}) {
	fmt.Printf("[debug]"+format+"\n", val...)
}
