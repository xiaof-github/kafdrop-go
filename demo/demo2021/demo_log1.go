package main

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)


func init() {
	path := "D:\\t.log"

	writer := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    200, // megabytes
		MaxBackups: 5,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}
	writers := []io.Writer{
		writer,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
}

func main() {
	log.Info("hello, world")
}