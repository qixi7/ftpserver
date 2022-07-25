package xlog

import (
	"encoding/json"
	"ftpServer/xcore/xlog/lumberjack-2.0"
	"log"
	"time"
)

type FileLogger struct {
	log    *log.Logger
	ticker *time.Ticker
	stop   chan struct{}
	writer *lumberjack.Logger
	queue  chan interface{}
}

func NewFileLogger(dir string, fname string) *FileLogger {
	writer := &lumberjack.Logger{
		Filename:   dir + "/" + fname,
		MaxSize:    100,
		MaxBackups: 100,
		MaxAge:     1,
	}

	log := log.New(writer, "", 0)
	ticker := time.NewTicker(time.Second * time.Duration(flushInterval))
	stop := make(chan struct{})
	queue := make(chan interface{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// 暂时do nothing
			case <-stop:
				return
			case item := <-queue:
				buf, err := json.Marshal(item)
				if err != nil {
					//Error(err)
					break
				}
				log.Output(0, string(buf))
			}
		}
	}()
	return &FileLogger{
		log:    log,
		ticker: ticker,
		stop:   stop,
		writer: writer,
		queue:  queue,
	}
}

func (fl *FileLogger) Close() {
	if fl == nil {
		return
	}
	fl.ticker.Stop()
	close(fl.stop)
	fl.writer.Close()
}

func (fl *FileLogger) WriteItem(item interface{}) {
	fl.queue <- item
}

func (fl *FileLogger) ReadItem(buf []byte, v interface{}) error {
	return json.Unmarshal(buf, v)
}
